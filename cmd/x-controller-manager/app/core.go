package app

import (
	"context"
	"fmt"
	"time"

	"github.com/fleezesd/xnightwatch/cmd/x-controller-manager/names"
	"k8s.io/client-go/metadata"
	"k8s.io/kubernetes/pkg/controller/garbagecollector"
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

// Start implement manager.Runnable interface.
func (gc *garbageCollector) Start(ctx context.Context) error {
	if _, err := startGarbageCollectorController(ctx, gc.cctx); err != nil {
		return err
	}
	return nil
}

func startGarbageCollectorController(ctx context.Context, controllerContext ControllerContext) (bool, error) {
	if !controllerContext.Config.ComponentConfig.GarbageCollectorController.EnableGarbageCollector {
		return false, nil
	}

	gcClientset := controllerContext.ClientBuilder.ClientOrDie("generic-garbage-collector")
	discoveryClient := controllerContext.ClientBuilder.DiscoveryClientOrDie("generic-garbage-collector")

	config := controllerContext.ClientBuilder.ConfigOrDie("generic-garbage-collector")
	// Increase garbage collector controller's throughput: each object deletion takes two API calls,
	// so to get |config.QPS| deletion rate we need to allow 2x more requests for this controller.
	config.QPS *= 2
	metadataClient, err := metadata.NewForConfig(config)
	if err != nil {
		return true, err
	}

	garbageCollector, err := garbagecollector.NewComposedGarbageCollector(
		ctx,
		gcClientset,
		metadataClient,
		controllerContext.RESTMapper,
		controllerContext.GraphBuilder,
	)
	if err != nil {
		return true, fmt.Errorf("failed to start the generic garbage collector: %w", err)
	}

	// Start the garbage collector.
	workers := int(controllerContext.Config.ComponentConfig.GarbageCollectorController.ConcurrentGCSyncs)
	const syncPeriod = 30 * time.Second
	go garbageCollector.Run(ctx, workers, syncPeriod)

	// Periodically refresh the RESTMapper with new discovery information and sync
	// the garbage collector.
	go garbageCollector.Sync(ctx, discoveryClient, syncPeriod)

	return true, nil
}
