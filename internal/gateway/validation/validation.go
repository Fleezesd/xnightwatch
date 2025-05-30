package validation

import (
	"context"

	"github.com/fleezesd/xnightwatch/internal/gateway/store"
	v1 "github.com/fleezesd/xnightwatch/pkg/api/gateway/v1"
	"github.com/google/wire"
)

// ProviderSet is a set of validator providers, used for dependency injection.
var ProviderSet = wire.NewSet(New, wire.Bind(new(any), new(*validator)))

// validator is a struct that implements the validate.IValidator interface.
type validator struct {
	ds store.IStore
}

// New is a factory function that creates and initializes the custom validator.
// It takes a store.IStore instance as input and returns *validator.
func New(ds store.IStore) (*validator, error) {
	vd := &validator{ds: ds}

	return vd, nil
}

// ValidateListMinerSetRequest is a method that validates the ListMinerSetRequest input.
func (vd *validator) ValidateListMinerSetRequest(ctx context.Context, req *v1.ListMinerSetRequest) error {
	return nil
}
