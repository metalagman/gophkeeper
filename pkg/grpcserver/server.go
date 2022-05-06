package grpcserver

import (
	"fmt"
	"google.golang.org/grpc"
	"gophkeeper/pkg/logger"
	"net"
)

type Server struct {
	config Config
	logger *logger.Logger
	server *grpc.Server
}

type Config struct {
	ListenAddr string `mapstructure:"listen_addr"`

	ServerOptions []grpc.ServerOption
}

type ServiceInit func(grpc.ServiceRegistrar)

func New(cfg Config) *Server {
	s := &Server{
		config: cfg,
		logger: logger.Global(),
	}
	s.server = grpc.NewServer(cfg.ServerOptions...)
	return s
}

func (s *Server) InitServices(init ...ServiceInit) {
	for _, f := range init {
		f(s.server)
	}
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.config.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	go func() {
		if err := s.server.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			s.logger.Fatal().Err(err).Send()
		}
	}()

	return nil
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}
