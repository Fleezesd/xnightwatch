package service

import (
	"context"

	v1 "github.com/fleezesd/xnightwatch/pkg/api/gateway/v1"
)

func (s *GatewayService) ListMinerSet(ctx context.Context, req *v1.ListMinerSetRequest) (*v1.ListMinerSetResponse, error) {
	mss, err := s.biz.MinerSets().List(ctx, "", req)
	if err != nil {
		return &v1.ListMinerSetResponse{}, err
	}

	return mss, nil
}
