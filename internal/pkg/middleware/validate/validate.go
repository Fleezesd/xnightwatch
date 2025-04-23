package validate

import (
	"context"

	"github.com/fleezesd/xnightwatch/pkg/api/zerrors"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
)

type validator interface {
	Validate() error
}

type IValidator interface {
	Validate(ctx context.Context, req any) error
}

// Validator is a validator middleware.
func Validator(vd IValidator) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			if v, ok := req.(validator); ok {
				// Kratos validation method
				if err := v.Validate(); err != nil {
					if se := new(errors.Error); errors.As(err, &se) {
						return nil, se
					}

					return nil, zerrors.ErrorInvalidParameter(err.Error()).WithCause(err)
				}
			}

			// Custom validation method
			if err := vd.Validate(ctx, req); err != nil {
				if se := new(errors.Error); errors.As(err, &se) {
					return nil, se
				}

				return nil, zerrors.ErrorInvalidParameter(err.Error()).WithCause(err)
			}

			return handler(ctx, req)
		}
	}
}
