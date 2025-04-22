package idempotent

import (
	"context"
	"errors"

	"github.com/fleezesd/xnightwatch/pkg/idempotent"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport"
)

func Idempotent(idt *idempotent.Idempotent) middleware.Middleware {
	return selector.Server(
		func(handler middleware.Handler) middleware.Handler {
			return func(ctx context.Context, rq any) (rp any, err error) {
				if tr, ok := transport.FromServerContext(ctx); ok {
					token := tr.RequestHeader().Get("X-Idempotent-ID")
					if token != "" {
						if idt.Check(ctx, token) {
							return handler(ctx, rq)
						}
						return nil, errors.New("idempotent token is invalid")
					}
				}

				return nil, errors.New("idempotent token is missing")
			}
		},
	).Match(idempotentBlacklist()).Build()
}

func idempotentBlacklist() selector.MatchFunc {
	return func(ctx context.Context, operation string) bool {
		return false
	}
}
