package options

import (
	"math"

	"github.com/fleezesd/xnightwatch/internal/nightwatch"
	kubeutil "github.com/fleezesd/xnightwatch/internal/pkg/util/kube"
	"github.com/fleezesd/xnightwatch/pkg/app"
	"github.com/fleezesd/xnightwatch/pkg/feature"
	"github.com/fleezesd/xnightwatch/pkg/log"
	genericoptions "github.com/fleezesd/xnightwatch/pkg/options"
	"github.com/fleezesd/xnightwatch/pkg/watch"
	"github.com/spf13/viper"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/tools/clientcmd"
	cliflag "k8s.io/component-base/cli/flag"
)

const (
	UserAgent = "xnightwatch"
)

var _ app.CliOptions = (*Options)(nil)

type Options struct {
	HealthOptions         *genericoptions.HealthOptions  `json:"health" mapstructure:"health"`
	MySQLOptions          *genericoptions.MySQLOptions   `json:"mysql" mapstructure:"mysql"`
	RedisOptions          *genericoptions.RedisOptions   `json:"redis" mapstructure:"redis"`
	MetricsOptions        *genericoptions.MetricsOptions `json:"metrics" mapstructure:"metrics"`
	WatchOptions          *watch.Options                 `json:"nightwatch" mapstructure:"nightwatch"`
	KubeConfig            string                         `json:"kubeconfig" mapstructure:"kubeconfig"`
	UserWatcherMaxWorkers int64                          `json:"user-watcher-max-workers"  mapstructure:"user-watcher-max-workers"`
	FeatureGates          map[string]bool                `json:"feature-gates" mapstructure:"feature-gates"`
	Log                   *log.Options                   `json:"log" mapstructure:"log"`
}

func NewOptions() *Options {
	o := &Options{
		HealthOptions:         genericoptions.NewHealthOptions(),
		MySQLOptions:          genericoptions.NewMySQLOptions(),
		RedisOptions:          genericoptions.NewRedisOptions(),
		MetricsOptions:        genericoptions.NewMetricsOptions(),
		WatchOptions:          watch.NewOptions(),
		UserWatcherMaxWorkers: math.MaxInt64,
		Log:                   log.NewOptions(),
	}

	return o
}

// Flags returns flags for a specific server by section name.
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	o.HealthOptions.AddFlags(fss.FlagSet("health"))
	o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))
	o.MetricsOptions.AddFlags(fss.FlagSet("metrics"))
	o.WatchOptions.AddFlags(fss.FlagSet("nightwatch"))
	o.Log.AddFlags(fss.FlagSet("log"))

	fs := fss.FlagSet("misc")
	fs.StringVar(&o.KubeConfig, "kubeconfig", o.KubeConfig, "Path to kubeconfig file with authorization and master location information.")
	fs.Int64Var(&o.UserWatcherMaxWorkers, "user-watcher-max-workers", o.UserWatcherMaxWorkers, "Specify the maximum concurrency event of user watcher.")
	feature.DefaultMutableFeatureGate.AddFlag(fs)

	return fss
}

func (o *Options) Complete() error {
	if err := viper.Unmarshal(&o); err != nil {
		return err
	}
	if o.UserWatcherMaxWorkers < 1 {
		o.UserWatcherMaxWorkers = math.MaxInt64
	}

	_ = feature.DefaultMutableFeatureGate.SetFromMap(o.FeatureGates)
	return nil
}

func (o *Options) Validate() error {
	errs := []error{}

	errs = append(errs, o.HealthOptions.Validate()...)
	errs = append(errs, o.MySQLOptions.Validate()...)
	errs = append(errs, o.RedisOptions.Validate()...)
	errs = append(errs, o.MetricsOptions.Validate()...)
	errs = append(errs, o.WatchOptions.Validate()...)
	errs = append(errs, o.Log.Validate()...)

	return utilerrors.NewAggregate(errs)
}

func (o *Options) ApplyTo(c *nightwatch.Config) error {
	c.MySQLOptions = o.MySQLOptions
	c.RedisOptions = o.RedisOptions
	c.WatchOptions = o.WatchOptions
	c.UserWatcherMaxWorkers = o.UserWatcherMaxWorkers
	return nil
}

// Config return xnightwatch config object.
func (o *Options) Config() (*nightwatch.Config, error) {
	// make clientset
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", o.KubeConfig)
	if err != nil {
		return nil, err
	}
	kubeutil.SetDefaultClientOptions(kubeutil.AddUserAgent(kubeconfig, UserAgent))
	if err != nil {
		return nil, err
	}

	c := &nightwatch.Config{}

	if err := o.ApplyTo(c); err != nil {
		return nil, err
	}
	return c, nil
}
