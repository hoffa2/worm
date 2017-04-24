package segment

import (
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/hoffa2/worm/chord"
)

type Remote struct {
	clientConns map[string]rpc.ChordClient
}

func (r Remote) tryDial(node *chord.Node) (chord.ChordClient, error) {

	dialerOpt := grpc.WithDialer(func(addr string, t time.Duration) (net.Conn, error) {
		return net.DialTimeout("udp", addr, t)
	})

	cc, err := grpc.Dial(node.GetRpcPort(), dialerOpt)
	if err != nil {
		return nil, err
	}

	return rpc.NewChordClient(cc), nil
}

func (r Remote) retrieveOrInitConn(node *chord.Node) (chord.ChordClient, error) {
	if conn, ok := r.clientConns[node.GetID()]; ok {
		return conn, nil
	}

	return r.tryDial(node)
}

func (r Remote) GetConn(node *rpc.Node) (chord.ChordClient, error) {
	return r.retrieveOrInitConn(node)
}
