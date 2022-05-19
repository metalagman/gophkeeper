package grpcservice

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "gophkeeper/api/proto"
	"gophkeeper/internal/server/model"
	storagemock "gophkeeper/internal/server/storage/mock"
	tokenmock "gophkeeper/pkg/token/mock"
	"net"
	"testing"
)

func TestIntegration(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uidOk := uuid.New()
	userOk := &model.User{
		ID:       uidOk,
		Email:    "user1@example.org",
		Password: "pass",
	}

	// mocking token manager
	tm := tokenmock.NewMockManager(ctrl)
	tm.EXPECT().Issue(userOk, gomock.Any()).AnyTimes().Return("token1", nil)

	// mocking user repo
	u := storagemock.NewMockUserRepository(ctrl)
	u.EXPECT().Create(gomock.Any(), &model.User{
		Email:    "user1@example.org",
		Password: "pass",
	}).AnyTimes().Return(
		&model.User{
			ID:       uidOk,
			Email:    "user1@example.org",
			Password: "pass",
		},
		nil,
	)

	// run mocked server
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}

	svc := NewAuth(u, tm)

	srv := grpc.NewServer()
	pb.RegisterAuthServer(srv, svc)
	go func() {
		if err := srv.Serve(l); err != nil && err != grpc.ErrServerStopped {
			panic(err) // We're in a goroutine - we can't t.Fatal/t.Error.
		}
	}()
	defer srv.GracefulStop()

	// real client for mocked service
	conn, err := grpc.Dial(l.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	cl := pb.NewAuthClient(conn)

	resp, err := cl.Register(ctx, &pb.RegisterRequest{
		Email:    "user1@example.org",
		Password: "pass",
	})
	assert.Equal(t, resp.Token, "token1")
}
