package biz

import (
	"github.com/fleezesd/xnightwatch/internal/gateway/biz/minerset"
	"github.com/fleezesd/xnightwatch/internal/gateway/store"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewBiz, wire.Bind(new(IBiz), new(*biz)))

type IBiz interface {
	MinerSets() minerset.MinerSetBiz
}

type biz struct {
	ds store.IStore
}

func NewBiz(ds store.IStore) *biz {
	return &biz{ds: ds}
}

func (b *biz) MinerSets() minerset.MinerSetBiz {
	return minerset.New(b.ds)
}
