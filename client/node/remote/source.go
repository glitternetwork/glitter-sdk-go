package remote

import (
	"context"
	"github.com/glitternetwork/glitter-sdk-go/client/node"

	"google.golang.org/grpc"
)

var (
	_ node.Source = &Source{}
)

// Source implements the keeper.Source interface relying on a GRPC connection
type Source struct {
	Ctx      context.Context
	GrpcConn *grpc.ClientConn
}

// NewSource returns a new Source instance
func NewSource(config *GRPCConfig) (*Source, error) {
	return &Source{
		Ctx:      context.Background(),
		GrpcConn: MustCreateGrpcConnection(config),
	}, nil
}

// Type implements keeper.Type
func (k Source) Type() string {
	return node.RemoteKeeper
}
