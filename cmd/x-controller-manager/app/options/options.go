package options

import (
	controllermanagerconfig "github.com/fleezesd/xnightwatch/cmd/x-controller-manager/app/config"
	"github.com/fleezesd/xnightwatch/cmd/x-controller-manager/names"
	ctrlmgrconfig "github.com/fleezesd/xnightwatch/internal/controller/apis/config"
	"github.com/fleezesd/xnightwatch/internal/controller/apis/config/latest"
	clientcmdutil "github.com/fleezesd/xnightwatch/internal/pkg/util/clientcmd"
	kubeutil "github.com/fleezesd/xnightwatch/internal/pkg/util/kube"
	clientset "github.com/fleezesd/xnightwatch/pkg/generated/clientset/versioned"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/client-go/tools/clientcmd"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/logs"
	logsapi "k8s.io/component-base/logs/api/v1"
	"k8s.io/component-base/metrics"
	"k8s.io/kubernetes/pkg/controller/garbagecollector"
)

const (
	// ControllerManagerUserAgent is the userAgent name when starting x-controller managers.
	ControllerManagerUserAgent = "x-controller-manager"
)

type Options struct {
	Generic                    *GenericControllerManagerConfigurationOptions
	GarbageCollectorController *GarbageCollectorControllerOptions
	ChainController            *ChainControllerOptions

	// ConfigFile is the location of the miner controller server's configuration file.
	ConfigFile string
	// WriteConfigTo is the path where the default configuration will be written.
	WriteConfigTo string
	// The address of the Kubernetes API server (overrides any value in kubeconfig).
	Master string
	// Path to kubeconfig file with authorization and master location information.
	Kubeconfig string
	Metrics    *metrics.Options
	Logs       *logs.Options
}

func NewOptions() (*Options, error) {
	componentConfig, err := latest.Default()
	if err != nil {
		return nil, err
	}

	o := Options{
		Generic:                    NewGenericControllerManagerConfigurationOptions(&componentConfig.Generic),
		GarbageCollectorController: NewGarbageCollectorControllerOptions(&componentConfig.GarbageCollectorController),
		ChainController:            NewChainControllerOptions(&componentConfig.ChainController),
		Kubeconfig:                 clientcmdutil.DefaultKubeconfig(),
		Metrics:                    metrics.NewOptions(),
		Logs:                       logs.NewOptions(),
	}

	gcIgnoredResources := make([]ctrlmgrconfig.GroupResource, 0, len(garbagecollector.DefaultIgnoredResources()))
	for r := range garbagecollector.DefaultIgnoredResources() {
		gcIgnoredResources = append(gcIgnoredResources, ctrlmgrconfig.GroupResource{Group: r.Group, Resource: r.Resource})
	}
	o.GarbageCollectorController.GCIgnoredResources = gcIgnoredResources
	o.Generic.LeaderElection.ResourceName = "x-controller-manager"
	o.Generic.LeaderElection.ResourceNamespace = "kube-system"

	return &o, nil
}

func (o *Options) Complete() error {
	return nil
}

// Flags returns flags for a specific APIServer by section name.
func (o *Options) Flags(allControllers []string, disabledControllers []string, controllerAliases map[string]string) cliflag.NamedFlagSets {
	fss := cliflag.NamedFlagSets{}
	o.Generic.AddFlags(&fss, allControllers, disabledControllers, controllerAliases)
	o.ChainController.AddFlags(fss.FlagSet(names.ChainController))
	o.GarbageCollectorController.AddFlags(fss.FlagSet(names.GarbageCollectorController))
	o.Metrics.AddFlags(fss.FlagSet("metrics"))
	logsapi.AddFlags(o.Logs, fss.FlagSet("logs"))

	fs := fss.FlagSet("misc")
	fs.StringVar(&o.ConfigFile, "config", o.ConfigFile, "The path to the configuration file.")
	fs.StringVar(&o.WriteConfigTo, "write-config-to", o.WriteConfigTo, "If set, write the default configuration values to this file and exit.")
	fs.StringVar(&o.Master, "master", o.Master, "The address of the Kubernetes API server (overrides any value in kubeconfig).")
	fs.StringVar(&o.Kubeconfig, "kubeconfig", o.Kubeconfig, "Path to kubeconfig file with authorization and master location information.")

	utilfeature.DefaultMutableFeatureGate.AddFlag(fss.FlagSet("generic"))
	return fss
}

// ApplyTo fills up x-controller manager config with options
func (o *Options) ApplyTo(c *controllermanagerconfig.Config, allControllers []string, disabledControllers []string, controllerAliases map[string]string) error {
	if err := o.Generic.ApplyTo(&c.ComponentConfig.Generic, allControllers, disabledControllers, controllerAliases); err != nil {
		return err
	}

	if err := o.GarbageCollectorController.ApplyTo(&c.ComponentConfig.GarbageCollectorController); err != nil {
		return err
	}

	if err := o.ChainController.ApplyTo(&c.ComponentConfig.ChainController); err != nil {
		return err
	}

	o.Metrics.Apply()

	return nil
}

// Validate is used to validaate the options and config before launching the controller.
func (o *Options) Validate(allControllers []string, disabledControllers []string, controllerAliases map[string]string) error {
	var errs []error

	errs = append(errs, o.Generic.Validate(allControllers, disabledControllers, controllerAliases)...)
	errs = append(errs, o.GarbageCollectorController.Validate()...)
	errs = append(errs, o.ChainController.Validate()...)

	return utilerrors.NewAggregate(errs)
}

// Config return a controller manager config objective
func (o Options) Config(allControllers []string, disabledControllers []string, controllerAliases map[string]string) (*controllermanagerconfig.Config, error) {
	kubeconfig, err := clientcmd.BuildConfigFromFlags(o.Master, o.Kubeconfig)
	if err != nil {
		return nil, err
	}
	kubeconfig.DisableCompression = true

	restConfig := kubeutil.AddUserAgent(kubeconfig, ControllerManagerUserAgent)
	client, err := clientset.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	c := &controllermanagerconfig.Config{
		KubeConfig: kubeutil.SetDefaultClientOptions(restConfig),
		Client:     client,
	}

	if err := o.ApplyTo(c, allControllers, disabledControllers, controllerAliases); err != nil {
		return nil, err
	}

	return c, nil
}
