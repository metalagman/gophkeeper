package grpcservice

import (
	"context"
	"google.golang.org/grpc"
	pb "gophkeeper/api/proto"
	"gophkeeper/internal/server/storage"
	"gophkeeper/pkg/grpcserver"
)

type Auth struct {
	pb.UnimplementedAuthServer

	users storage.UserRepository
}

func NewAuth(u storage.UserRepository) *Auth {
	return &Auth{
		users: u,
	}
}

func (s Auth) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s Auth) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Auth) Init() grpcserver.ServiceInit {
	return func(registrar grpc.ServiceRegistrar) {
		pb.RegisterAuthServer(registrar, s)
	}
}
