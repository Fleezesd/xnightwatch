package options

import (
	"github.com/fleezesd/xnightwatch/internal/gateway"
	"github.com/fleezesd/xnightwatch/pkg/app"
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
	JaegerOptions *genericoptions.JaegerOptions  `json:"jaeger" mapstructure:"jaeger"`
	Metrics       *genericoptions.MetricsOptions `json:"metrics" mapstructure:"metrics"`
	EnableTLS     bool                           `json:"enable-tls" mapstructure:"enable-tls"`
	Kubeconfig    string                         `json:"kubeconfig" mapstructure:"kubeconfig"`
	Log           *log.Options                   `json:"log" mapstructure:"log"`
}

func NewOptions() *Options {
	o := &Options{}

	return o
}

// Flags returns flags for a specific server by section name.
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	return fss
}

// Compute completes all the required options.
func (o *Options) Complete() error {
	return nil
}

// Validate checks Options and return a slice of found errs.
func (o *Options) Validate() error {
	errs := []error{}

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
