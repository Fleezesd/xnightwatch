package store

import (
	"context"

	"github.com/fleezesd/xnightwatch/internal/gateway/model"
	"github.com/fleezesd/xnightwatch/internal/pkg/meta"
)

// ChainStore defines the chain storage interface.
type ChainStore interface {
	Create(ctx context.Context, ch *model.ChainM) error
	Get(ctx context.Context, filters map[string]any) error
	Update(ctx context.Context, ch *model.ChainM) error
	Delete(ctx context.Context, filters map[string]any) error
	List(ctx context.Context, namespace string, opts ...meta.ListOption) (int64, []*model.ChainM, error)
}

// chainStore is a structure which implements the ChainStore interface.
type chainStore struct {
	ds *datastore
}

// newChainStore creates a new chainStore instance with provided datastore.
func newChainStore(ds *datastore) *chainStore {
	return &chainStore{ds: ds}
}

// db is an alias for m.ds.Core(ctx context.Context), a convenience method to get the core database instance.
