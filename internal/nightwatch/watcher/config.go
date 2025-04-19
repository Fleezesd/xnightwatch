package watcher

import "github.com/fleezesd/xnightwatch/internal/pkg/client/store"

// AggregateConfig aggregates the configurations of all watchers and serves as a configuration aggregator.
type AggregatedConfig struct {
	// The purpose of nightwatch is to handle asynchronous tasks on the platform
	// in a unified manner, so a store aggregation type is needed here.
	Store store.Interface

	// Then maximum concurrency event of user watcher.
	UserWatcherMaxWorkers int64
}
