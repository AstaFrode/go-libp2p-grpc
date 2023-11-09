package libp2pgrpc

import (
	"context"

	gostream "github.com/AstaFrode/go-libp2p-gostream"
	"github.com/AstaFrode/go-libp2p/core/host"
	"google.golang.org/grpc"
)

// ServerOption allows for functional setting of options on a Server.
type ServerOption func(*Server)

var _ grpc.ServiceRegistrar = &Server{}

type Server struct {
	host host.Host
	grpc *grpc.Server
	ctx  context.Context
}

// NewGrpcServer creates a Server object with the given LibP2P host
// and protocol.
func NewGrpcServer(ctx context.Context, h host.Host, opts ...ServerOption) (*Server, error) {
	grpcServer := grpc.NewServer()

	srv := &Server{
		host: h,
		ctx:  ctx,
		grpc: grpcServer,
	}

	for _, opt := range opts {
		opt(srv)
	}

	listener, err := gostream.Listen(srv.host, ProtocolID)
	if err != nil {
		return nil, err
	}

	go srv.grpc.Serve(listener)

	return srv, nil
}

func (s *Server) RegisterService(serviceDesc *grpc.ServiceDesc, srv interface{}) {
	s.grpc.RegisterService(serviceDesc, srv)
}
