package service

import (
	"context"

	v1 "github.com/fleezesd/xnightwatch/pkg/api/gateway/v1"
	"github.com/fleezesd/xnightwatch/pkg/version"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *GatewayService) GetVersion(ctx context.Context, rq *emptypb.Empty) (*v1.GetVersionResponse, error) {
	vinfo := version.Get()
	return &v1.GetVersionResponse{
		GitVersion:   vinfo.GitVersion,
		GitCommit:    vinfo.GitCommit,
		GitTreeState: vinfo.GitTreeState,
		BuildDate:    vinfo.BuildDate,
		GoVersion:    vinfo.GoVersion,
		Compiler:     vinfo.Compiler,
		Platform:     vinfo.Platform,
	}, nil
}
