package options

import (
	cmconfig "github.com/fleezesd/xnightwatch/internal/controller/apis/config"
	cliflag "k8s.io/component-base/cli/flag"
)

// GenericControllerManagerConfigurationOptions holds the options which are generic.
type GenericControllerManagerConfigurationOptions struct {
	*cmconfig.GenericControllerManagerConfiguration
}

func NewGenericControllerManagerConfigurationOptions(cfg *cmconfig.GenericControllerManagerConfiguration) *GenericControllerManagerConfigurationOptions {
	return &GenericControllerManagerConfigurationOptions{
		GenericControllerManagerConfiguration: cfg,
	}
}

// AddFlags adds flags related to ChainController for controller manager to the specified FlagSet.
func (o *GenericControllerManagerConfigurationOptions) AddFlags(
	fss *cliflag.NamedFlagSets,
	allControllers []string,
	disabledControllers []string,
	controllerAliasesmap map[string]string,
) {

}
