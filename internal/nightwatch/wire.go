//go:build wireinject
// +build wireinject

package nightwatch

import (
	"github.com/fleezesd/xnightwatch/internal/pkg/client/store"
	"github.com/fleezesd/xnightwatch/pkg/db"
	"github.com/google/wire"

	gwstore "github.com/fleezesd/xnightwatch/internal/gateway/store"
	ucstore "github.com/fleezesd/xnightwatch/internal/usercenter/store"
)

//go:generate go run github.com/google/wire/cmd/wire

func wireStoreClient(*db.MySQLOptions) (store.Interface, error) {
	wire.Build(
		db.ProviderSet,
		store.ProviderSet,
		gwstore.ProviderSet,
		ucstore.ProviderSet,
	)

	return nil, nil
}
