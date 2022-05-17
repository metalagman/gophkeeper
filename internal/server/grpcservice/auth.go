package grpcservice

import (
	"context"
	"google.golang.org/grpc"
	pb "gophkeeper/api/proto"
	"gophkeeper/internal/server/model"
	"gophkeeper/internal/server/storage"
	"gophkeeper/pkg/grpcserver"
	"gophkeeper/pkg/token"
	"time"
)

const tokenLifetime = time.Hour * 24 * 365

type Auth struct {
	pb.UnimplementedAuthServer

	users storage.UserRepository
	token token.Manager
}

func NewAuth(u storage.UserRepository, tm token.Manager) *Auth {
	return &Auth{
		users: u,
		token: tm,
	}
}

func (s Auth) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	u := &model.User{
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
	}

	u, err := s.users.Create(ctx, u)
	if err != nil {
		return nil, err
	}

	t, err := s.token.Issue(u, tokenLifetime)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		Token: t,
	}, nil
}

func (s Auth) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	u, err := s.users.ReadByEmailAndPassword(ctx, request.GetEmail(), request.GetPassword())
	if err != nil {
		return nil, err
	}

	t, err := s.token.Issue(u, tokenLifetime)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		Token: t,
	}, nil
}

func (s *Auth) Init() grpcserver.ServiceInit {
	return func(registrar grpc.ServiceRegistrar) {
		pb.RegisterAuthServer(registrar, s)
	}
}
