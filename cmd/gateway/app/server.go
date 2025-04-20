package app

import (
	"github.com/fleezesd/xnightwatch/cmd/gateway/app/options"
	"github.com/fleezesd/xnightwatch/internal/gateway"
	"github.com/fleezesd/xnightwatch/pkg/app"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

const commandDesc = `The gateway server is the back-end portal server. All 
requests from the front-end will arrive at the gateway, requests will be uniformly processed 
and distributed by the gateway.`

// NewApp creates an App object with default parameters.
func NewApp() *app.App {
	opts := options.NewOptions()
	application := app.NewApp(gateway.Name, "Launch a gateway server",
		app.WithDescription(commandDesc),
		app.WithOptions(opts),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
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
func Run(c *gateway.Config, stopCh <-chan struct{}) error {
	server, err := c.Complete().New(stopCh)
	if err != nil {
		return err
	}

	return server.Run(stopCh)
}
