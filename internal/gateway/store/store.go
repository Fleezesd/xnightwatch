package store

//go:generate mockgen -self_package=github.com/fleezesd/xnightwatch/internal/gateway/store -destination=mock_store.go -package=store github.com/fleezesd/xnightwatch/internal/gateway/store IStore,ChainStore,MinerStore,MinerSetStore

import (
	"context"
	"sync"

	"github.com/google/wire"
	"gorm.io/gorm"
)

var ProviderSet = wire.NewSet(NewStore, wire.Bind(new(IStore), new(*datastore)))

var (
	once sync.Once
	S    *datastore
)

type transactionKey struct{}

// IStore is an interface defining the required methods for a Store.
type IStore interface {
	TX(context.Context, func(ctx context.Context) error) error
	Chains() ChainStore
	Miners() MinerStore
	MinerSets() minerSetStore
}

// datastore is a concrete implementation of IStore interface.
type datastore struct {
	// core is the main database, use the name `core` to indicate that this is the main database.
	core *gorm.DB
}

// Verify that datastore implements IStore interface.
var _ IStore = (*datastore)(nil)

// NewStore initializes a new datastore by using the given gorm.DB and returns it.
func NewStore(db *gorm.DB) *datastore {
	once.Do(func() {
		S = &datastore{db}
	})

	return S
}

// Core returns the core gorm.DB from the datastore. If there is an ongoing transaction,
// the transaction's gorm.DB is returned instead.
func (ds *datastore) Core(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(transactionKey{}).(*gorm.DB)
	if ok {
		return tx
	}

	return ds.core
}

// TX is a method to execute a function inside a transaction, it takes a context and a function as parameters.
func (ds *datastore) TX(ctx context.Context, fn func(ctx context.Context) error) error {
	return ds.core.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			ctx = context.WithValue(ctx, transactionKey{}, tx)
			return fn(ctx)
		},
	)
}

// Chains returns a ChainStore that interacts with datastore.
func (ds *datastore) Chains() ChainStore {
	return newChainStore(ds)
}

// Miners returns a MinerStore that interacts with datastore.
func (ds *datastore) Miners() MinerStore {
	return newMinerStore(ds)
}

// MinerSets returns a MinerSetStore that interacts with datastore.
func (ds *datastore) MinerSets() minerSetStore {
	return *newMinerSetStore(ds)
}
