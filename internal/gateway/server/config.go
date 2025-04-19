package server

import (
	genericoptions "github.com/fleezesd/xnightwatch/pkg/options"
)

type Config struct {
	HTTP genericoptions.HTTPOptions
	GRPC genericoptions.GRPCOptions
}
