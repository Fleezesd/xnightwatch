package service

import (
	"github.com/fleezesd/xnightwatch/internal/gateway/biz"
	v1 "github.com/fleezesd/xnightwatch/pkg/api/gateway/v1"
)

type GatewayService struct {
	v1.UnimplementedGatewayServer
	biz biz.IBiz
}

func NewGatewayService(biz biz.IBiz) *GatewayService {
	return &GatewayService{
		biz: biz,
	}
}
