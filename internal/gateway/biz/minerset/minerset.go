package minerset

import (
	"context"

	"github.com/fleezesd/xnightwatch/internal/gateway/store"
	"github.com/fleezesd/xnightwatch/internal/pkg/meta"
	v1 "github.com/fleezesd/xnightwatch/pkg/api/gateway/v1"
	"github.com/fleezesd/xnightwatch/pkg/log"
	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// MinerSetBiz defines functions used to handle minerset rquest.
type MinerSetBiz interface {
	List(ctx context.Context, namespace string, req *v1.ListMinerSetRequest) (*v1.ListMinerSetResponse, error)
}

type minerSetBiz struct {
	ds store.IStore
}

var _ MinerSetBiz = (*minerSetBiz)(nil)

func New(ds store.IStore) *minerSetBiz {
	return &minerSetBiz{ds}
}

func (b *minerSetBiz) List(ctx context.Context, namespace string, req *v1.ListMinerSetRequest) (*v1.ListMinerSetResponse, error) {
	total, list, err := b.ds.MinerSets().List(ctx, namespace, meta.WithOffset(req.Offset), meta.WithLimit(req.Limit))
	if err != nil {
		log.C(ctx).Errorw(err, "Failed to list minerset")
		return nil, err
	}

	mss := make([]*v1.MinerSet, 0, len(list))
	for _, item := range list {
		var ms v1.MinerSet
		_ = copier.Copy(&ms, &item)
		ms.CreatedAt = timestamppb.New(item.CreatedAt)
		ms.UpdatedAt = timestamppb.New(item.UpdatedAt)
		mss = append(mss, &ms)
	}

	return &v1.ListMinerSetResponse{
		TotalCount: total,
		MinerSets:  mss,
	}, nil
}
