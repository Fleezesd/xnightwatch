package nightwatch

import (
	"github.com/fleezesd/xnightwatch/internal/nightwatch/watcher"
	"github.com/fleezesd/xnightwatch/pkg/db"
	"github.com/fleezesd/xnightwatch/pkg/log"
	genericoptions "github.com/fleezesd/xnightwatch/pkg/options"
	"github.com/fleezesd/xnightwatch/pkg/watch"
	"github.com/fleezesd/xnightwatch/pkg/watch/logger/x"
	"github.com/jinzhu/copier"
	"k8s.io/apimachinery/pkg/util/wait"

	// trigger init functions in `internal/nightwatch/watcher/all`
	_ "github.com/fleezesd/xnightwatch/internal/nightwatch/watcher/all"
)

type nightWatch struct {
	// watch
	*watch.Watch
}

// Config is the configuration for the nightwatch server.
type Config struct {
	MySQLOptions *genericoptions.MySQLOptions
	RedisOptions *genericoptions.RedisOptions
	WatchOptions *watch.Options
	// The maximum concurrency event of user watcher.
	UserWatcherMaxWorkers int64
}

// CompletedConfig same as Config, just to swap private object.
type CompletedConfig struct {
	*Config
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (c *Config) Complete() CompletedConfig {
	return CompletedConfig{c}
}

func (c *Config) CreateWatcherConfig() (*watcher.AggregatedConfig, error) {
	var mysqlOptions db.MySQLOptions
	_ = copier.Copy(&mysqlOptions, c.MySQLOptions)
	storeClient, err := wireStoreClient(&mysqlOptions)
	if err != nil {
		log.Errorw(err, "Failed to create MySQL client")
		return nil, err
	}

	return &watcher.AggregatedConfig{
		Store:                 storeClient,
		UserWatcherMaxWorkers: c.UserWatcherMaxWorkers,
	}, nil
}

func (c *Config) New() (*nightWatch, error) {
	client, err := c.RedisOptions.NewClient()
	if err != nil {
		log.Errorw(err, "Failed to create Redis client")
		return nil, err
	}

	cfg, err := c.CreateWatcherConfig()
	if err != nil {
		return nil, err
	}

	initalize := watcher.NewWatcherInitializer(cfg.Store, cfg.UserWatcherMaxWorkers)
	opts := []watch.Option{
		watch.WithInitialize(initalize),
		watch.WithLogger(x.NewLogger()),
		watch.WithLockName("x-nightwatch-lock"),
	}

	nw, err := watch.NewWatch(c.WatchOptions, client, opts...)
	if err != nil {
		return nil, err
	}

	return &nightWatch{
		nw,
	}, nil
}

// Run keep retrying to acquire lock and then start the Cron job.
func (nw *nightWatch) Run(stopCh <-chan struct{}) {
	nw.Start(wait.ContextForChannel(stopCh))
	// graceful shutdown
	<-stopCh

	nw.Stop()
}
