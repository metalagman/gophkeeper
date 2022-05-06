package app

import (
	"fmt"
	"gophkeeper/internal/server/config"
	"gophkeeper/pkg/grpcserver"
	"gophkeeper/pkg/logger"
)

type App struct {
	config config.Config
	logger logger.Logger
	stop   chan struct{}
	server *grpcserver.Server
}

func New(cfg config.Config) (*App, error) {
	l := *logger.Global()

	s := grpcserver.New(cfg.GRPC)
	if err := s.Start(); err != nil {
		return nil, fmt.Errorf("grpc: %w", err)
	}

	a := &App{
		config: cfg,
		logger: l,
		stop:   make(chan struct{}),
		server: s,
	}

	return a, nil
}

func (a *App) Stop() {
	close(a.stop)
	a.server.Stop()
}
