package auth

import (
	"context"

	"github.com/fleezesd/xnightwatch/internal/gateway/locales"
	jwtutil "github.com/fleezesd/xnightwatch/internal/pkg/util/jwt"
	"github.com/fleezesd/xnightwatch/pkg/api/zerrors"
	"github.com/fleezesd/xnightwatch/pkg/i18n"
	"github.com/fleezesd/xnightwatch/pkg/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

type AuthProvider interface {
	Auth(ctx context.Context, token string, obj, act string) (userID string, allowed bool, err error)
}

// Auth is a authentication and authorization middleware.
func Auth(a AuthProvider) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (resp any, err error) {
			accessToken := jwtutil.TokenFromServerContext(ctx)
			if tr, ok := transport.FromServerContext(ctx); ok {
				_, allowed, err := a.Auth(ctx, accessToken, "*", tr.Operation())
				if err != nil {
					log.Errorw(err, "Authorization failure occurs", "operation", tr.Operation())
					return nil, err
				}
				if !allowed {
					return nil, zerrors.ErrorForbidden(i18n.FromContext(ctx).T(locales.NoPermission))
				}
			}
			return handler(ctx, req)
		}
	}
}
