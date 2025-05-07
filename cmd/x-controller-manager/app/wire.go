//go:build wireinject
// +build wireinject

package app

//go:generate go run github.com/google/wire/cmd/wire
import (
	"github.com/fleezesd/xnightwatch/internal/gateway/store"
	"github.com/fleezesd/xnightwatch/pkg/db"
	"github.com/google/wire"
)

func wireStoreClient(mysqlOptions *db.MySQLOptions) (store.IStore, error) {
	wire.Build(
		db.ProviderSet,
		store.ProviderSet,
	)

	return nil, nil
}
