package store

import (
	"context"
	"errors"

	"github.com/fleezesd/xnightwatch/internal/gateway/model"
	"github.com/fleezesd/xnightwatch/internal/pkg/meta"
	"gorm.io/gorm"
)

// ChainStore defines the chain storage interface.
type ChainStore interface {
	Create(ctx context.Context, ch *model.ChainM) error
	Get(ctx context.Context, filters map[string]any) (*model.ChainM, error)
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
func (d *chainStore) db(ctx context.Context) *gorm.DB {
	return d.ds.Core(ctx)
}

// Create creates a new chain record in the database.
func (d *chainStore) Create(ctx context.Context, ch *model.ChainM) error {
	return d.db(ctx).Create(ch).Error
}

// Get retrieves a single chain record from the database based on provided filters.
func (d *chainStore) Get(ctx context.Context, filters map[string]any) (*model.ChainM, error) {
	chain := &model.ChainM{}
	if err := d.db(ctx).Where(filters).First(chain).Error; err != nil {
		return nil, err
	}
	return chain, nil
}

// Update updates a chain record in the database.
func (d *chainStore) Update(ctx context.Context, ch *model.ChainM) error {
	return d.db(ctx).Save(ch).Error
}

// Delete deletes a chain record from the database based on provided filters.
func (d *chainStore) Delete(ctx context.Context, filters map[string]any) error {
	err := d.db(ctx).Where(filters).Delete(&model.ChainM{}).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return nil
}

// List returns a list of chain records according to the provided query conditions.
func (d *chainStore) List(ctx context.Context, namespace string, opts ...meta.ListOption) (count int64, ret []*model.ChainM, err error) {
	los := meta.NewListOptions(opts...)
	if namespace != "" {
		los.Filters["namespace"] = namespace
	}
	ans := d.db(ctx).
		Where(los.Filters).
		Offset(los.Offset).
		Limit(los.Limit).
		Order("id desc").
		Find(&ret).
		Offset(-1).
		Limit(-1).
		Count(&count)

	return count, ret, ans.Error
}
