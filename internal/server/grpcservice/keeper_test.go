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
	"gophkeeper/internal/server/storage"
	storagemock "gophkeeper/internal/server/storage/mock"
	"gophkeeper/pkg/grpcserver"
	"gophkeeper/pkg/usercontext"
	"log"
	"testing"
)

var (
	okUserID = uuid.New()
)

func TestIntegrationKeeper_Create(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cl, stop := getTestClient(t, ctrl)
	defer stop()

	resp, err := cl.CreateSecret(ctx, &pb.CreateSecretRequest{
		Name:    "secret1",
		Type:    "raw",
		Content: []byte("keepitsecret"),
	})
	assert.NoError(t, err)
	assert.Equal(t, resp.Name, "secret1")

	t.Log("Done integration testing")
}

func TestIntegrationKeeper_ReadByName(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cl, stop := getTestClient(t, ctrl)
	defer stop()

	resp, err := cl.ReadSecret(ctx, &pb.ReadSecretRequest{
		Name: "secret1",
	})
	assert.NoError(t, err)
	assert.Equal(t, resp.Name, "secret1")
	assert.Equal(t, resp.Type, "raw")

	t.Log("Done integration testing")
}

func TestIntegrationKeeper_DeleteByName(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cl, stop := getTestClient(t, ctrl)
	defer stop()

	_, err := cl.DeleteSecret(ctx, &pb.DeleteSecretRequest{
		Name: "secret1",
	})
	assert.NoError(t, err)

	t.Log("Done integration testing")
}

func TestIntegrationKeeper_List(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cl, stop := getTestClient(t, ctrl)
	defer stop()

	resp, err := cl.ListSecrets(ctx, &pb.ListSecretsRequest{})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.GetSecrets())
	assert.Equal(t, resp.GetSecrets()[0].Name, "secret1")
	assert.Equal(t, resp.GetSecrets()[0].Type, "raw")

	t.Log("Done integration testing")
}

func getTestClient(t *testing.T, ctrl *gomock.Controller) (pb.KeeperClient, func()) {
	secrets := getTestSecretRepository(ctrl)
	svc := NewKeeper(secrets)

	s := grpcserver.New(
		grpcserver.WithListenAddr("localhost:0"),
		grpcserver.WithServices(svc),
		grpcserver.WithAuthFunc(testAuthFunc),
	)
	if err := s.Start(); err != nil {
		t.Fatal(err)
	}

	t.Log("Server started")

	// real client for mocked service
	conn, err := grpc.Dial(s.ListenAddr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}

	stop := func() {
		_ = conn.Close()
		s.Stop()
	}

	cl := pb.NewKeeperClient(conn)
	return cl, stop
}

func getTestSecretRepository(ctrl *gomock.Controller) storage.SecretRepository {
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
	secrets.EXPECT().ReadByName(gomock.Any(), okUserID, "secret1").AnyTimes().Return(&model.Secret{
		Name:    "secret1",
		Type:    "raw",
		Content: []byte("keepitsecret"),
	}, nil)
	secrets.EXPECT().DeleteByName(gomock.Any(), okUserID, "secret1").AnyTimes().Return(nil)
	secrets.EXPECT().List(gomock.Any(), okUserID).AnyTimes().Return([]*model.Secret{
		&model.Secret{
			Name: "secret1",
			Type: "raw",
		},
	}, nil)

	return secrets
}

func testAuthFunc(ctx context.Context) (context.Context, error) {
	log.Println("test auth func")
	ctx = usercontext.WriteUID(ctx, okUserID)
	return ctx, nil
}
