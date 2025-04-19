package clean

import (
	"github.com/fleezesd/xnightwatch/internal/nightwatch/watcher"
	"github.com/fleezesd/xnightwatch/internal/pkg/client/store"
	"github.com/fleezesd/xnightwatch/pkg/watch"
)

var _ watch.Watcher = (*cleanWatcher)(nil)

type cleanWatcher struct {
	store store.Interface
}

// Run runs the watcher.
func (w *cleanWatcher) Run() {

}

func (w *cleanWatcher) SetAggregateConfig(config *watcher.AggregatedConfig) {
	w.store = config.Store
}

func init() {
	watch.Register("clean", &cleanWatcher{})
}
