package store

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
	Chains() ChainStore
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

// Chains returns a ChainStore that interacts with datastore.
func (ds *datastore) Chains() ChainStore {
	return nil
}
