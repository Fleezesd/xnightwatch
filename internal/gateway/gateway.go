package gateway

import (
	"github.com/fleezesd/xnightwatch/pkg/log"
	"github.com/go-kratos/kratos/v2"
)

var (
	Name = "gateway"
)

type Config struct {
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (cfg *Config) Complete() completedConfig {
	return completedConfig{cfg}
}

type completedConfig struct {
	*Config
}

func (c completedConfig) New(stopCh <-chan struct{}) (*Server, error) {
	return &Server{}, nil
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
