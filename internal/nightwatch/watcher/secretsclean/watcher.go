package secretsclean

import (
	"github.com/fleezesd/xnightwatch/internal/nightwatch/watcher"
	"github.com/fleezesd/xnightwatch/internal/pkg/client/store"
	"github.com/fleezesd/xnightwatch/pkg/watch"
)

var _ watch.Watcher = (*secretsCleanWatcher)(nil)

// watcher implement.
type secretsCleanWatcher struct {
	store store.Interface
}

func (w *secretsCleanWatcher) Run() {

}

func (w *secretsCleanWatcher) SetAggregateConifg(config *watcher.AggregatedConfig) {
	w.store = config.Store
}

func init() {
	watch.Register("secretsclean", &secretsCleanWatcher{})
}
