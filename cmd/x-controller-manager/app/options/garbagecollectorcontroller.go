package options

import (
	garbagecollectorconfig "github.com/fleezesd/xnightwatch/internal/controller/apis/config"
	"github.com/spf13/pflag"
)

// GarbageCollectorControllerOptions holds the GarbageCollectorController options.
type GarbageCollectorControllerOptions struct {
	*garbagecollectorconfig.GarbageCollectorControllerConfiguration
}

func NewGarbageCollectorControllerOptions(cfg *garbagecollectorconfig.GarbageCollectorControllerConfiguration) *GarbageCollectorControllerOptions {
	return &GarbageCollectorControllerOptions{
		GarbageCollectorControllerConfiguration: cfg,
	}
}

// AddFlags adds flags related to GarbageCollectorController for controller manager to the specified FlagSet.
func (o *GarbageCollectorControllerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.Int32Var(&o.ConcurrentGCSyncs, "concurrent-gc-syncs", o.ConcurrentGCSyncs, "The number of garbage collector workers that are allowed to sync concurrently.")
	fs.BoolVar(&o.EnableGarbageCollector, "enable-garbage-collector", o.EnableGarbageCollector, "Enables the generic garbage collector. MUST be synced with the corresponding flag of the kube-apiserver.")
}

// ApplyTo fills up GarbageCollectorController config with options.
func (o *GarbageCollectorControllerOptions) ApplyTo(cfg *garbagecollectorconfig.GarbageCollectorControllerConfiguration) error {
	if o == nil {
		return nil
	}

	cfg.ConcurrentGCSyncs = o.ConcurrentGCSyncs
	cfg.GCIgnoredResources = o.GCIgnoredResources
	cfg.EnableGarbageCollector = o.EnableGarbageCollector

	return nil
}

// Validate checks validation of GarbageCollectorController.
func (o *GarbageCollectorControllerOptions) Validate() []error {
	if o == nil {
		return nil
	}

	errs := []error{}
	return errs
}
