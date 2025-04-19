package service

import (
	v1 "github.com/fleezesd/xnightwatch/pkg/api/gateway/v1"
)

type GatewayService struct {
	v1.UnimplementedGatewayServer
}

func NewGatewayService() *GatewayService {
	return &GatewayService{}
}
