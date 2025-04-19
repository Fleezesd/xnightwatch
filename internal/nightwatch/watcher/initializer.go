package watcher

import (
	"github.com/fleezesd/xnightwatch/internal/pkg/client/store"
	"github.com/fleezesd/xnightwatch/pkg/watch"
)

// WatcherInitializer is used for initialization of the specific watcher plugins.
type WatcherInitializer struct {
	// The purpose of nightwatch is to handle asynchronous tasks on the onex platform
	// in a unified manner, so a store aggregation type is needed here.
	store store.Interface

	// Then maximum concurrency event of user watcher.
	UserWatcherMaxWorkers int64
}

var _ watch.WatcherInitializer = &WatcherInitializer{}

func NewWatcherInitializer(store store.Interface, userWatcherMaxWorkers int64) *WatcherInitializer {
	return &WatcherInitializer{
		store:                 store,
		UserWatcherMaxWorkers: userWatcherMaxWorkers,
	}
}

func (w *WatcherInitializer) Initialize(wc watch.Watcher) {
	if wants, ok := wc.(WantsStore); ok {
		wants.SetStore(w.store)
	}

	if wants, ok := wc.(WantsAggregateConfig); ok {
		wants.SetAggregateConifg(&AggregatedConfig{
			Store:                 w.store,
			UserWatcherMaxWorkers: w.UserWatcherMaxWorkers,
		})
	}
}
