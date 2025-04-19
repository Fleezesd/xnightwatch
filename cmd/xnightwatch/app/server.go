package app

import (
	"github.com/fleezesd/xnightwatch/cmd/xnightwatch/app/options"
	"github.com/fleezesd/xnightwatch/internal/nightwatch"
	"github.com/fleezesd/xnightwatch/pkg/app"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

const commandDesc = `The nightwatch server is responsible for executing some async tasks 
like linux cronjob. You can add Cron(github.com/robfig/cron) jobs on the given schedule
use the Cron spec format.`

func NewApp() *app.App {
	opts := options.NewOptions()

	application := app.NewApp("xnightwatch", "Launch a x asynchronous task processing server",
		app.WithDescription(commandDesc),
		app.WithOptions(opts),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
		app.WithHealthCheckFunc(func() error {
			go opts.HealthOptions.ServeHealthCheck()
			return nil
		}),
	)
	return application
}

func run(opts *options.Options) app.RunFunc {
	return func() error {
		cfg, err := opts.Config()
		if err != nil {
			return err
		}
		return Run(cfg, genericapiserver.SetupSignalHandler())
	}
}

// Run runs the specified APIServer. This should never exit.
func Run(c *nightwatch.Config, stopCh <-chan struct{}) error {
	nw, err := c.Complete().New()
	if err != nil {
		return err
	}

	nw.Run(stopCh)
	return nil
}
