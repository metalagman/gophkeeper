package grpcservice

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "gophkeeper/api/proto"
	"gophkeeper/internal/server/model"
	"gophkeeper/internal/server/storage"
	"gophkeeper/pkg/apperr"
	"gophkeeper/pkg/token"
	"time"
)

const tokenLifetime = time.Hour * 24 * 365

type User struct {
	pb.UnimplementedUserServer

	users storage.UserRepository
	token token.Manager
}

func NewUser(u storage.UserRepository, tm token.Manager) *User {
	return &User{
		users: u,
		token: tm,
	}
}

func (s User) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	u := &model.User{
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
	}

	u, err := s.users.Create(ctx, u)
	if err != nil {
		return nil, err
	}

	t, err := s.token.Issue(u, tokenLifetime)
	switch err {
	case nil:
		// all is ok
	case apperr.ErrConflict:
		return nil, status.Error(codes.InvalidArgument, err.Error())
	default:
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RegisterResponse{
		Token: t,
	}, nil
}

func (s User) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	u, err := s.users.ReadByEmailAndPassword(ctx, request.GetEmail(), request.GetPassword())
	if err != nil {
		return nil, err
	}

	t, err := s.token.Issue(u, tokenLifetime)
	switch err {
	case nil:
		// all is ok
	case apperr.ErrNotFound:
		return nil, status.Error(codes.Unauthenticated, err.Error())
	default:
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.LoginResponse{
		Token: t,
	}, nil
}

func (s *User) RegisterService(r grpc.ServiceRegistrar) {
	pb.RegisterUserServer(r, s)
}

func (s *User) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return ctx, nil
}
