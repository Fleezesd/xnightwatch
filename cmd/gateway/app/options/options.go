package options

import (
	"github.com/fleezesd/xnightwatch/internal/gateway"
	"github.com/fleezesd/xnightwatch/pkg/app"
	"github.com/fleezesd/xnightwatch/pkg/client"
	"github.com/fleezesd/xnightwatch/pkg/feature"
	"github.com/fleezesd/xnightwatch/pkg/log"
	genericoptions "github.com/fleezesd/xnightwatch/pkg/options"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	cliflag "k8s.io/component-base/cli/flag"
)

const (
	UserAgent = "gateway"
)

var _ app.CliOptions = (*Options)(nil)

type Options struct {
	GRPCOptions   *genericoptions.GRPCOptions    `json:"grpc" mapstructure:"grpc"`
	HTTPOptions   *genericoptions.HTTPOptions    `json:"http" mapstructure:"http"`
	TLSOptions    *genericoptions.TLSOptions     `json:"tls" mapstructure:"tls"`
	MySQLOptions  *genericoptions.MySQLOptions   `json:"mysql" mapstructure:"mysql"`
	RedisOptions  *genericoptions.RedisOptions   `json:"redis" mapstructure:"redis"`
	EtcdOptions   *genericoptions.EtcdOptions    `json:"etcd" mapstructure:"etcd"`
	JaegerOptions *genericoptions.JaegerOptions  `json:"jaeger" mapstructure:"jaeger"`
	ConsulOptions *genericoptions.ConsulOptions  `json:"consul" mapstructure:"consul"`
	Metrics       *genericoptions.MetricsOptions `json:"metrics" mapstructure:"metrics"`
	EnableTLS     bool                           `json:"enable-tls" mapstructure:"enable-tls"`
	Kubeconfig    string                         `json:"kubeconfig" mapstructure:"kubeconfig"`
	FeatureGates  map[string]bool                `json:"feature-gates" mapstructure:"feature-gates"`
	Log           *log.Options                   `json:"log" mapstructure:"log"`
}

func NewOptions() *Options {
	o := &Options{
		GRPCOptions:   genericoptions.NewGRPCOptions(),
		HTTPOptions:   genericoptions.NewHTTPOptions(),
		TLSOptions:    genericoptions.NewTLSOptions(),
		MySQLOptions:  genericoptions.NewMySQLOptions(),
		RedisOptions:  genericoptions.NewRedisOptions(),
		EtcdOptions:   genericoptions.NewEtcdOptions(),
		JaegerOptions: genericoptions.NewJaegerOptions(),
		ConsulOptions: genericoptions.NewConsulOptions(),
		Metrics:       genericoptions.NewMetricsOptions(),
		Log:           log.NewOptions(),
	}

	return o
}

// Flags returns flags for a specific server by section name.
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	o.GRPCOptions.AddFlags(fss.FlagSet("grpc"))
	o.HTTPOptions.AddFlags(fss.FlagSet("http"))
	o.TLSOptions.AddFlags(fss.FlagSet("tls"))
	o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))
	o.EtcdOptions.AddFlags(fss.FlagSet("etcd"))
	o.JaegerOptions.AddFlags(fss.FlagSet("jaeger"))
	o.ConsulOptions.AddFlags(fss.FlagSet("consul"))
	o.Metrics.AddFlags(fss.FlagSet("metrics"))
	o.Log.AddFlags(fss.FlagSet("log"))

	fs := fss.FlagSet("misc")
	client.AddFlags(fs)
	fs.StringVar(&o.Kubeconfig, "kubeconfig", o.Kubeconfig, "Path to kubeconfig file with authorization and master location information.")
	feature.DefaultMutableFeatureGate.AddFlag(fs)
	return fss
}

// Compute completes all the required options.
func (o *Options) Complete() error {
	if o.JaegerOptions.ServiceName == "" {
		o.JaegerOptions.ServiceName = UserAgent
	}
	// set feature gates from options
	_ = feature.DefaultMutableFeatureGate.SetFromMap(o.FeatureGates)
	return nil
}

// Validate checks Options and return a slice of found errs.
func (o *Options) Validate() error {
	errs := []error{}

	errs = append(errs, o.GRPCOptions.Validate()...)
	errs = append(errs, o.HTTPOptions.Validate()...)
	errs = append(errs, o.TLSOptions.Validate()...)
	errs = append(errs, o.MySQLOptions.Validate()...)
	errs = append(errs, o.RedisOptions.Validate()...)
	errs = append(errs, o.EtcdOptions.Validate()...)
	errs = append(errs, o.JaegerOptions.Validate()...)
	errs = append(errs, o.ConsulOptions.Validate()...)
	errs = append(errs, o.Metrics.Validate()...)
	errs = append(errs, o.Log.Validate()...)

	return utilerrors.NewAggregate(errs)
}

// ApplyTo fills up gateway config with options
func (o *Options) ApplyTo(c *gateway.Config) error {
	return nil
}

// Config returns gateway config object
func (o *Options) Config() (*gateway.Config, error) {

	c := &gateway.Config{}

	if err := o.ApplyTo(c); err != nil {
		return nil, err
	}
	return c, nil
}
