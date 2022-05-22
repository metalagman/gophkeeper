package grpcservice

import (
	"google.golang.org/grpc"
	pb "gophkeeper/api/proto"
)

type Keeper struct {
	pb.UnimplementedKeeperServer
}

func NewKeeper() *Keeper {
	return &Keeper{}
}

func (s *Keeper) RegisterService(r grpc.ServiceRegistrar) {
	pb.RegisterKeeperServer(r, s)
}
