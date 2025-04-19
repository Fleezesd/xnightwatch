package store

import (
	"sync"

	"github.com/google/wire"
	"gorm.io/gorm"
)

var ProviderSet = wire.NewSet(NewStore, wire.Bind(new(IStore), new(*datastore)))

var (
	once sync.Once
	S    *datastore
)

// IStore is an interface defining the required methods for a Store.
type IStore interface {
}

// datastore is a concrete implementation of IStore interface.
type datastore struct {
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
