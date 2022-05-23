package grpcservice

import (
	"context"
	"github.com/google/uuid"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gophkeeper/pkg/token"
	"gophkeeper/pkg/usercontext"
)

func BuildAuthFunc(tok token.Manager) grpc_auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		mdt, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, err
		}

		uid, err := tok.Decode(mdt)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
		}

		u, err := uuid.Parse(uid.Identity())
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
		}

		ctx = usercontext.WriteUID(ctx, u)
		return ctx, nil
	}
}
