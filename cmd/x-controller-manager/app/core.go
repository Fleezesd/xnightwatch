package app

import (
	"context"

	"github.com/fleezesd/xnightwatch/cmd/x-controller-manager/names"
	ctrl "sigs.k8s.io/controller-runtime"
)

func newGarbageCollectorControllerDescriptor() *ControllerDescriptor {
	return &ControllerDescriptor{
		name:    names.GarbageCollectorController,
		aliases: []string{"garbagecollector"},
	}
}

func addGarbageCollectorController(ctx context.Context, mgr ctrl.Manager, cctx ControllerContext) (bool, error) {
	return true, mgr.Add(&garbageCollector{cctx: cctx})
}

type garbageCollector struct {
	cctx ControllerContext
}

func (gc *garbageCollector) Start(ctx context.Context) error {
	return nil
}

func startGarbageCollectorController(ctx context.Context, cctx ControllerContext) (bool, error) {
	return true, nil
}
