package watch

type WatcherInitializer interface {
	Initialize(watcher Watcher)
}
