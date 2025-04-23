package gateway

import (
	"os"

	"github.com/fleezesd/xnightwatch/internal/gateway/server"
	"github.com/fleezesd/xnightwatch/internal/pkg/bootstrap"
	"github.com/fleezesd/xnightwatch/pkg/db"
	"github.com/fleezesd/xnightwatch/pkg/log"
	"github.com/fleezesd/xnightwatch/pkg/version"
	"github.com/go-kratos/kratos/v2"
	"github.com/jinzhu/copier"
	"k8s.io/client-go/rest"

	genericoptions "github.com/fleezesd/xnightwatch/pkg/options"
)

var (
	// Name is the name of the compiled software.
	Name = "gateway"

	ID, _ = os.Hostname()
)

type Config struct {
	GRPCOptions   *genericoptions.GRPCOptions
	HTTPOptions   *genericoptions.HTTPOptions
	TLSOptions    *genericoptions.TLSOptions
	MySQLOptions  *genericoptions.MySQLOptions
	RedisOptions  *genericoptions.RedisOptions
	EtcdOptions   *genericoptions.EtcdOptions
	JaegerOptions *genericoptions.JaegerOptions
	ConsulOptions *genericoptions.ConsulOptions

	KubeConfig *rest.Config
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (cfg *Config) Complete() completedConfig {
	return completedConfig{cfg}
}

type completedConfig struct {
	*Config
}

func (c completedConfig) New(stopCh <-chan struct{}) (*Server, error) {
	if err := c.JaegerOptions.SetTracerProvider(); err != nil {
		return nil, err
	}

	appInfo := bootstrap.NewAppInfo(ID, Name, version.Get().String())
	conf := &server.Config{
		HTTP: *c.HTTPOptions,
		GRPC: *c.GRPCOptions,
		TLS:  *c.TLSOptions,
	}

	var mysqlOptions db.MySQLOptions
	var redisOptions db.RedisOptions
	_ = copier.Copy(&mysqlOptions, c.MySQLOptions)
	_ = copier.Copy(&redisOptions, c.RedisOptions)

	app, cleanup, err := wireApp(stopCh, appInfo, conf, &mysqlOptions, &redisOptions, c.EtcdOptions)
	if err != nil {
		return nil, err
	}
	defer cleanup()
	return &Server{
		app: app,
	}, nil
}

// Server represents the gateway server.
type Server struct {
	app *kratos.App
}

func (s *Server) Run(stopCh <-chan struct{}) error {
	go func() {
		if err := s.app.Run(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-stopCh

	log.Infof("Gracefully shutting down server ...")

	if err := s.app.Stop(); err != nil {
		log.Errorw(err, "Failed to gracefully shutdown kratos application")
		return err
	}
	return nil
}
