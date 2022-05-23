package grpcservice

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "gophkeeper/api/proto"
	"gophkeeper/internal/server/model"
	"gophkeeper/internal/server/storage"
	"gophkeeper/pkg/apperr"
	"gophkeeper/pkg/usercontext"
)

type Keeper struct {
	pb.UnimplementedKeeperServer

	secrets storage.SecretRepository
}

func NewKeeper(s storage.SecretRepository) *Keeper {
	return &Keeper{
		secrets: s,
	}
}

func (s *Keeper) RegisterService(r grpc.ServiceRegistrar) {
	pb.RegisterKeeperServer(r, s)
}

func (s *Keeper) CreateSecret(ctx context.Context, request *pb.CreateSecretRequest) (*pb.CreateSecretResponse, error) {
	uid := usercontext.ReadUID(ctx)
	if !uid.Valid {
		return nil, status.Error(codes.Unauthenticated, apperr.ErrUnauthorized.Error())
	}

	m := &model.Secret{
		UserID:  uid.UUID,
		Name:    request.GetName(),
		Type:    request.GetType(),
		Content: request.GetContent(),
	}
	if m, err := s.secrets.Create(ctx, uid.UUID, m); err != nil {
		if errors.Is(err, apperr.ErrConflict) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	} else {
		return &pb.CreateSecretResponse{
			Name: m.Name,
			Type: m.Type,
		}, nil
	}
}

func (s *Keeper) ReadSecret(ctx context.Context, request *pb.ReadSecretRequest) (*pb.ReadSecretResponse, error) {
	uid := usercontext.ReadUID(ctx)
	if !uid.Valid {
		return nil, status.Error(codes.Unauthenticated, apperr.ErrUnauthorized.Error())
	}

	if m, err := s.secrets.ReadByName(ctx, uid.UUID, request.GetName()); err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	} else {
		return &pb.ReadSecretResponse{
			Name:    m.Name,
			Type:    m.Type,
			Content: m.Content,
		}, nil
	}
}

func (s *Keeper) DeleteSecret(ctx context.Context, request *pb.DeleteSecretRequest) (*pb.DeleteSecretResponse, error) {
	uid := usercontext.ReadUID(ctx)
	if !uid.Valid {
		return nil, status.Error(codes.Unauthenticated, apperr.ErrUnauthorized.Error())
	}

	if err := s.secrets.DeleteByName(ctx, uid.UUID, request.GetName()); err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	} else {
		return &pb.DeleteSecretResponse{}, nil
	}
}

func (s *Keeper) ListSecrets(ctx context.Context, request *pb.ListSecretsRequest) (*pb.ListSecretsResponse, error) {
	uid := usercontext.ReadUID(ctx)
	if !uid.Valid {
		return nil, status.Error(codes.Unauthenticated, apperr.ErrUnauthorized.Error())
	}

	mm, err := s.secrets.List(ctx, uid.UUID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &pb.ListSecretsResponse{}
	for _, m := range mm {
		resp.Secrets = append(resp.Secrets, &pb.SecretDescription{
			Name: m.Name,
			Type: m.Type,
		})
	}

	return resp, nil
}
