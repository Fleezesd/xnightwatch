package watcher

import (
	"github.com/fleezesd/xnightwatch/internal/nightwatch/watcher"
	"github.com/fleezesd/xnightwatch/internal/pkg/client/store"
	"github.com/fleezesd/xnightwatch/pkg/watch"
)

// watcher implement.
type userWatcher struct {
	store      store.Interface
	maxWorkers int64
}

func (w *userWatcher) Run() {
}

func (w *userWatcher) SetAggregateConifg(config *watcher.AggregatedConfig) {
	w.store = config.Store
	w.maxWorkers = config.UserWatcherMaxWorkers
}

func init() {
	watch.Register("user", &userWatcher{})
}
