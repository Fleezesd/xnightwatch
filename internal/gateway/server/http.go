package server

import (
	"github.com/fleezesd/xnightwatch/internal/gateway/service"
	"github.com/fleezesd/xnightwatch/internal/pkg/pprof"
	v1 "github.com/fleezesd/xnightwatch/pkg/api/gateway/v1"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *Config, gw *service.GatewayService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.HTTP.Network != "" {
		opts = append(opts, http.Network(c.HTTP.Network))
	}
	if c.HTTP.Addr != "" {
		opts = append(opts, http.Address(c.HTTP.Addr))
	}
	if c.HTTP.Timeout != 0 {
		opts = append(opts, http.Timeout(c.HTTP.Timeout))
	}
	srv := http.NewServer(opts...)
	h := openapiv2.NewHandler()
	srv.HandlePrefix("/openapi", h)
	srv.Handle("/metrcis", promhttp.Handler())
	srv.Handle("", pprof.NewHandler())

	v1.RegisterGatewayHTTPServer(srv, gw)
	return srv
}
