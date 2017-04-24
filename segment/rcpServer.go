package segment

import (
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
func SetupRPCServer(chord chord.ChordServer, port string) (*RpcServer, error) {
	conn, err := net.Listen("udp", "0.0.0.0:"+port)
	if err != nil {
		return nil, err
	}

	server := grpc.NewServer()
	chord.RegisterChordServer(server, chord)

	go func() {
		server.Serve(conn)
	}()

	goss := &RpcServer{
		rpcServer: server,
	}

	return goss, nil
}

func (r RpcServer) CloseServer() {
	r.rpcServer.GracefulStop()
}
