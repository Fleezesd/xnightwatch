package config

import (
	ctrlmgrconfig "github.com/fleezesd/xnightwatch/internal/controller/apis/config"
)

// Config is the main context object for the controller
type Config struct {
	ComponentConfig *ctrlmgrconfig.XControllerManagerConfiguration
}

// CompletedConfig same as Config, just to swap private object.
type CompletedConfig struct {
	*Config
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (c *Config) Complete() *CompletedConfig {
	return &CompletedConfig{c}
}
