package options

import (
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/logs"
)

const (
	// ControllerManagerUserAgent is the userAgent name when starting x-controller managers.
	ControllerManagerUserAgent = "x-controller-manager"
)

type Options struct {
	Logs *logs.Options
}

func NewOptions() (*Options, error) {
	return &Options{}, nil
}

func (o *Options) Complete() error {
	return nil
}

// Flags returns flags for a specific APIServer by section name.
func (o *Options) Flags() cliflag.NamedFlagSets {
	fss := cliflag.NamedFlagSets{}
	return fss
}
