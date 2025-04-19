package watch

import "github.com/spf13/pflag"

// Options contains everything necessary to create and run a watch server.
type Options struct {
	DisableWatchers []string `json:"disable-watchers" mapstructure:"disable-watchers"`
}

func NewOptions() *Options {
	o := &Options{
		DisableWatchers: []string{},
	}

	return o
}

// Flags returns flags for a specific server by section name.
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringSliceVar(&o.DisableWatchers, "disable-watchers", o.DisableWatchers, "The list of watchers that should be disabled.")
}

// Validate validates all the required options.
func (o *Options) Validate() []error {
	errs := []error{}

	return errs
}
