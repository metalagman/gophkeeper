package grpcservice

import (
	"context"
	"github.com/google/uuid"
	grpczerolog "github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2"
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/v2/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gophkeeper/pkg/logger"
	"gophkeeper/pkg/token"
	"gophkeeper/pkg/usercontext"
)

func BuildAuthFunc(tok token.Manager) grpcauth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		mdt, err := grpcauth.AuthFromMD(ctx, "bearer")
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

func BuildUnaryInterceptors() []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		logging.UnaryServerInterceptor(grpczerolog.InterceptorLogger(logger.Global().Logger)),
		recovery.UnaryServerInterceptor(),
	}
}
