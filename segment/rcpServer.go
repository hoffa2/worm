package segment

import (
	"context"
	"errors"
	"net"

	"github.com/hoffa2/worm/chord"

	"google.golang.org/grpc"
)

var (
	ErrAddrInvalid = errors.New("Could not resolve addr")
)

type RpcServer struct {
	rpcServer *grpc.Server
}

// SetupGossip start listening for aliens
// param port port on which to accept connections
func SetupRPCServer(chord chord.ChordServer, port string) (*Gossip, error) {
	conn, err := net.Listen("udp", "0.0.0.0:"+port)
	if err != nil {
		return nil, err
	}

	server := grpc.NewServer()
	chord.RegisterChordServer(server, chord)

	go func() {
		server.Serve(conn)
	}()

	goss := &Gossip{
		rpcServer: server,
	}

	return goss, nil
}

func (r RpcServer) CloseServer() {
	r.rpcServer.GracefulStop()
}

func (r RpcServer) Alive(ctx context.Context, empty *chord.Empty) (*chord.Alive, error) {

}

func (r RpcServer) FindPredecessor(ctx context.Context, node *chord.Node) (*chord.FromNode, error) {

}

func (r RpcServer) Init(ctx context.Context, node *chord.Node) (*chord.Empty, error) {

}

func (r RpcServer) Notify(ctx context.Context, node *chord.Node) (*chord.Empty, error) {

}

func (r RpcServer) Shutdown(ctx context.Context, empty *chord.Empty) (*chord.Empty, error) {

}
