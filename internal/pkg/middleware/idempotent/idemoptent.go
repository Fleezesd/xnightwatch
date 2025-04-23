package idempotent

import (
	"context"

	v1 "github.com/fleezesd/xnightwatch/pkg/api/gateway/v1"
	"github.com/fleezesd/xnightwatch/pkg/api/zerrors"
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
						return nil, zerrors.ErrorIdepotentTokenExpired("idempotent token is invalid")
					}
				}

				return nil, zerrors.ErrorIdempotentMissingToken("idempotent token is missing")
			}
		},
	).Match(idempotentBlacklist()).Build()
}

func idempotentBlacklist() selector.MatchFunc {
	blacklist := make(map[string]struct{})
	blacklist[v1.OperationGatewayCreateMiner] = struct{}{}
	blacklist[v1.OperationGatewayCreateMinerSet] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := blacklist[operation]; ok {
			return true
		}
		return false
	}
}
