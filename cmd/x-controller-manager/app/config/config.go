package config

import (
	restclient "k8s.io/client-go/rest"

	ctrlmgrconfig "github.com/fleezesd/xnightwatch/internal/controller/apis/config"
	clientset "github.com/fleezesd/xnightwatch/pkg/generated/clientset/versioned"
)

// Config is the main context object for the controller
type Config struct {
	ComponentConfig *ctrlmgrconfig.XControllerManagerConfiguration

	// the general clientset
	Client clientset.Interface

	// the rest config for the master
	KubeConfig *restclient.Config
}

// CompletedConfig same as Config, just to swap private object.
type CompletedConfig struct {
	*Config
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (c *Config) Complete() *CompletedConfig {
	return &CompletedConfig{c}
}
