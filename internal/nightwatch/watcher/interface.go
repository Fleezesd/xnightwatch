package watcher

import (
	"github.com/fleezesd/xnightwatch/internal/pkg/client/store"
	"github.com/fleezesd/xnightwatch/pkg/watch"
)

type WantsAggregateConfig interface {
	watch.Watcher
	SetAggregateConifg(config *AggregatedConfig)
}

type WantsStore interface {
	watch.Watcher
	SetStore(store store.Interface)
}
