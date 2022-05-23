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
	"gophkeeper/pkg/grpcserver"
	tokenmock "gophkeeper/pkg/token/mock"
	"gophkeeper/pkg/usercontext"
	"log"
	"testing"
)

var (
	okUserID = uuid.New()
)

func TestIntegrationKeeper(t *testing.T) {
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
	secrets := storagemock.NewMockSecretRepository(ctrl)
	secrets.EXPECT().Create(gomock.Any(), okUserID, &model.Secret{
		UserID:  okUserID,
		Name:    "secret1",
		Type:    "raw",
		Content: []byte("keepitsecret"),
	}).AnyTimes().Return(&model.Secret{
		UserID:  okUserID,
		Name:    "secret1",
		Type:    "raw",
		Content: []byte("keepitsecret"),
	}, nil)

	svc := NewKeeper(secrets)

	s := grpcserver.New(
		grpcserver.WithListenAddr("localhost:0"),
		grpcserver.WithServices(svc),
		grpcserver.WithAuthFunc(testAuthFunc),
	)
	if err := s.Start(); err != nil {
		t.Fatal(err)
	}

	log.Println("start ok")

	// real client for mocked service
	conn, err := grpc.Dial(s.ListenAddr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
		s.Stop()
	}(conn)

	cl := pb.NewKeeperClient(conn)

	resp, err := cl.CreateSecret(ctx, &pb.CreateSecretRequest{
		Name:    "secret1",
		Type:    "raw",
		Content: []byte("keepitsecret"),
	})
	assert.NoError(t, err)
	assert.Equal(t, resp.Name, "secret1")

	t.Log("done")
}

func testAuthFunc(ctx context.Context) (context.Context, error) {
	log.Println("test auth func")
	ctx = usercontext.WriteUID(ctx, okUserID)
	return ctx, nil
}
