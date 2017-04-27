package main

import (
	"errors"
	"net"

	"github.com/hoffa2/worm/protobuf/chord"

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
func SetupRPCServer(s chord.ChordServer, port string, l net.Listener) (*RpcServer, error) {
	server := grpc.NewServer()
	chord.RegisterChordServer(server, s)

	go func() {
		server.Serve(l)
	}()

	goss := &RpcServer{
		rpcServer: server,
	}

	return goss, nil
}

func (r RpcServer) CloseServer() {
	r.rpcServer.GracefulStop()
}
