package segment

import (
	"errors"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/hoffa2/worm/protobuf/chord"
)

var (
	ErrNotReachable = errors.New("Host is not reachable")
)

type Reachable interface {
	IsReachable(string) bool
}

type ClientRemote struct {
	clientConns map[string]chord.ChordClient
	Reachable
}

func SetupRemote(r Reachable) *ClientRemote {
	return &ClientRemote{
		Reachable:   r,
		clientConns: make(map[string]chord.ChordClient),
	}
}

func (c *ClientRemote) tryDial(node *chord.Node) (chord.ChordClient, error) {
	dialerOpt := grpc.WithDialer(func(addr string, t time.Duration) (net.Conn, error) {
		return net.DialTimeout("udp", addr, t)
	})

	cc, err := grpc.Dial(node.GetRpcPort(), dialerOpt)
	if err != nil {
		return nil, err
	}

	return chord.NewChordClient(cc), nil
}

func (c *ClientRemote) retrieveOrInitConn(node *chord.Node) (chord.ChordClient, error) {
	if conn, ok := c.clientConns[node.GetID()]; ok {
		return conn, nil
	}

	return c.tryDial(node)
}

func (c *ClientRemote) GetConn(node *chord.Node) (chord.ChordClient, error) {
	if !c.IsReachable(node.GetIpAddress()) {
		return nil, ErrNotReachable
	}
	return c.retrieveOrInitConn(node)
}
