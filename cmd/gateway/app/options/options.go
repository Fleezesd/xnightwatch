package options

import (
	"github.com/fleezesd/xnightwatch/internal/gateway"
	"github.com/fleezesd/xnightwatch/pkg/app"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	cliflag "k8s.io/component-base/cli/flag"
)

const (
	UserAgent = "gateway"
)

var _ app.CliOptions = (*Options)(nil)

type Options struct {
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
