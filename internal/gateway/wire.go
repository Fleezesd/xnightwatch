// go:build wireinject
//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package gateway

//go:generate go run github.com/google/wire/cmd/wire

import (
	"github.com/fleezesd/xnightwatch/internal/gateway/biz"
	"github.com/fleezesd/xnightwatch/internal/gateway/server"
	"github.com/fleezesd/xnightwatch/internal/gateway/service"
	"github.com/fleezesd/xnightwatch/internal/gateway/store"
	customvalidation "github.com/fleezesd/xnightwatch/internal/gateway/validation"
	"github.com/fleezesd/xnightwatch/internal/pkg/bootstrap"
	"github.com/fleezesd/xnightwatch/internal/pkg/idempotent"
	"github.com/fleezesd/xnightwatch/internal/pkg/validation"
	"github.com/fleezesd/xnightwatch/pkg/db"
	genericoptions "github.com/fleezesd/xnightwatch/pkg/options"
	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"
)

func wireApp(
	<-chan struct{},
	bootstrap.AppInfo,
	*server.Config,
	*db.MySQLOptions,
	*db.RedisOptions,
	*genericoptions.EtcdOptions,
) (*kratos.App, func(), error) {
	wire.Build(
		bootstrap.NewEtcdRegistrar,
		bootstrap.ProviderSet,
		server.ProviderSet,
		store.ProviderSet,
		db.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		validation.ProviderSet,
		idempotent.ProviderSet,
		customvalidation.ProviderSet,
	)

	return nil, nil, nil
}
