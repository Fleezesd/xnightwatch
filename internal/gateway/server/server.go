package server

import (
	"context"
	"encoding/json"

	"github.com/fleezesd/xnightwatch/pkg/idempotent"
	"github.com/fleezesd/xnightwatch/pkg/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"

	xmetrics "github.com/fleezesd/xnightwatch/internal/pkg/metrics"
	krtlog "github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel/metric"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewServers, NewGRPCServer, NewHTTPServer)

func NewServers(hs *http.Server, gs *grpc.Server) []transport.Server {
	return []transport.Server{hs, gs}
}

func NewMiddlewares(logger krtlog.Logger, idt *idempotent.Idempotent, meter metric.Meter) []middleware.Middleware {
	return []middleware.Middleware{
		recovery.Recovery(
			recovery.WithHandler(func(ctx context.Context, req, err any) error {
				data, _ := json.Marshal(req)
				log.C(ctx).Errorw(err.(error), "Catching a panic", "req", string(data))
				return nil
			}),
		),
		metrics.Server(
			metrics.WithSeconds(xmetrics.KratosMetricSeconds),
			metrics.WithRequests(xmetrics.KratosServerMetricRequests),
		),
	}
}
