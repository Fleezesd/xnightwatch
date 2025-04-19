package store

import (
	"sync"

	"github.com/google/wire"

	gwstore "github.com/fleezesd/xnightwatch/internal/gateway/store"
	ucstore "github.com/fleezesd/xnightwatch/internal/usercenter/store"
)

// ProviderSet is store providers.
var ProviderSet = wire.NewSet(NewStore, wire.Bind(new(Interface), new(*datastore)))

var (
	once sync.Once
	S    *datastore
)

type datastore struct {
	gw gwstore.IStore
	uc ucstore.IStore
}

var _ Interface = (*datastore)(nil)

// Interface defines the storage interface.
type Interface interface {
	Gateway() gwstore.IStore
	UserCenter() ucstore.IStore
}

func (ds *datastore) Gateway() gwstore.IStore {
	return ds.gw
}

func (ds *datastore) UserCenter() ucstore.IStore {
	return ds.uc
}

func NewStore(gw gwstore.IStore, uc ucstore.IStore) *datastore {
	once.Do(func() {
		S = &datastore{
			gw: gw,
			uc: uc,
		}
	})
	return S
}
