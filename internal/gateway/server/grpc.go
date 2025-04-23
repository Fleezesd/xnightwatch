package server

import (
	"github.com/fleezesd/xnightwatch/internal/gateway/service"
	v1 "github.com/fleezesd/xnightwatch/pkg/api/gateway/v1"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *Config, gw *service.GatewayService, middlewares []middleware.Middleware) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			middlewares...,
		),
	}
	if c.GRPC.Network != "" {
		opts = append(opts, grpc.Network(c.GRPC.Network))
	}
	if c.GRPC.Addr != "" {
		opts = append(opts, grpc.Address(c.GRPC.Addr))
	}
	if c.GRPC.Timeout != 0 {
		opts = append(opts, grpc.Timeout(c.GRPC.Timeout))
	}
	if c.TLS.UseTLS {
		opts = append(opts, grpc.TLSConfig(c.TLS.MustTLSConfig()))
	}

	srv := grpc.NewServer(opts...)
	v1.RegisterGatewayServer(srv, gw)
	return srv
}
