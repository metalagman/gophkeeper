package grpcservice

import (
	"google.golang.org/grpc"
	pb "gophkeeper/api/proto"
	"gophkeeper/internal/server/storage"
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
