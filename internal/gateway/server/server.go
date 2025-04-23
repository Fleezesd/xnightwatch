package server

import (
	"context"
	"encoding/json"

	authmw "github.com/fleezesd/xnightwatch/internal/pkg/middleware/auth"
	idempotentmw "github.com/fleezesd/xnightwatch/internal/pkg/middleware/idempotent"
	"github.com/fleezesd/xnightwatch/internal/pkg/middleware/logging"
	"github.com/fleezesd/xnightwatch/internal/pkg/middleware/validate"
	"github.com/fleezesd/xnightwatch/pkg/i18n"
	"github.com/fleezesd/xnightwatch/pkg/idempotent"
	"github.com/fleezesd/xnightwatch/pkg/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"golang.org/x/text/language"

	"github.com/fleezesd/xnightwatch/internal/gateway/locales"
	xmetrics "github.com/fleezesd/xnightwatch/internal/pkg/metrics"
	i18nmw "github.com/fleezesd/xnightwatch/internal/pkg/middleware/i18n"
	krtlog "github.com/go-kratos/kratos/v2/log"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewServers, NewGRPCServer, NewHTTPServer, NewMiddlewares)

func NewServers(hs *http.Server, gs *grpc.Server) []transport.Server {
	return []transport.Server{hs, gs}
}

func NewMiddlewares(logger krtlog.Logger, idt *idempotent.Idempotent, a authmw.AuthProvider, v validate.IValidator) []middleware.Middleware {
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
		i18nmw.Translator(i18n.WithLanguage(language.English), i18n.WithFS(locales.Locales)),
		idempotentmw.Idempotent(idt),
		ratelimit.Server(),
		tracing.Server(),
		selector.Server(authmw.Auth(a)).Match(NewWhiteListMatcher()).Build(),
		validate.Validator(v),
		logging.Server(logger),
	}
}

func NewWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})

	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}
