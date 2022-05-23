package grpcserver

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"gophkeeper/pkg/logger"
	"log"
	"net"
)

type Server struct {
	listenAddr string

	logger             *logger.Logger
	server             *grpc.Server
	services           []Service
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
	authFunc           grpc_auth.AuthFunc
}

func (s *Server) ListenAddr() string {
	return s.listenAddr
}

type ServerOption func(*Server)

func WithListenAddr(a string) ServerOption {
	return func(server *Server) {
		server.listenAddr = a
	}
}

func WithAuthFunc(af grpc_auth.AuthFunc) ServerOption {
	return func(server *Server) {
		server.authFunc = af
	}
}

func WithUnaryInterceptors(in ...grpc.UnaryServerInterceptor) ServerOption {
	return func(server *Server) {
		server.unaryInterceptors = append(server.unaryInterceptors, in...)
	}
}

func WithStreamInterceptors(in ...grpc.StreamServerInterceptor) ServerOption {
	return func(server *Server) {
		server.streamInterceptors = append(server.streamInterceptors, in...)
	}
}

func WithServices(in ...Service) ServerOption {
	return func(server *Server) {
		server.services = append(server.services, in...)
	}
}

type Service interface {
	RegisterService(grpc.ServiceRegistrar)
}

func New(opts ...ServerOption) *Server {
	s := &Server{
		logger: logger.Global(),
	}

	for _, o := range opts {
		o(s)
	}

	if s.authFunc != nil {
		s.unaryInterceptors = append(s.unaryInterceptors, grpc_auth.UnaryServerInterceptor(s.authFunc))
		s.streamInterceptors = append(s.streamInterceptors, grpc_auth.StreamServerInterceptor(s.authFunc))
	}

	return s
}

func (s *Server) RegisterServices(services ...Service) {
	for _, svc := range services {
		svc.RegisterService(s.server)
	}
}

func (s *Server) Start() error {
	log.Printf("%+v", s)
	lis, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	// required for testing
	s.listenAddr = lis.Addr().String()

	s.server = grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				s.unaryInterceptors...,
			),
		),
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				s.streamInterceptors...,
			),
		),
	)

	s.RegisterServices(s.services...)

	go func() {
		s.logger.Info().Str("host", s.listenAddr).Msg("Listening incoming GRPC connections")
		if err := s.server.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			s.logger.Fatal().Err(err).Send()
		}
	}()

	return nil
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}
