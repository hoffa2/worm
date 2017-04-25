package segment

import (
	"context"
	"io"
	"os"

	"crypto/sha1"

	"github.com/hoffa2/worm/protobuf/chord"
)

type Node interface {
	Shutdown()
}

type Node struct {
	*chord.Node
	Node
}

func hashValue(str string) string {
	h := sha1.New()
	io.WriteString(h, str)
	return string(h.Sum(nil))
}

func SetupKNode(addr, port string, nodeinterface Node) (*Kademlia, error) {
	h, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	n := &chord.Node{
		ID:        hashValue(h),
		IpAddress: addr,
		RpcPort:   port,
	}
	return &Kademlia{Node: n}
}

func (n *Node) Alive(ctx context.Context, empty *chord.Empty) (*chord.Alive, error) {
	return &chord.Alive{IsAlive: true}, nil
}

func (n *Node) FindPredecessor(ctx context.Context, empty *chord.Empty) (*chord.Alive, error) {
	return nil, nil
}

func (n *Node) Init(ctx context.Context, empty *chord.Empty) (*chord.Alive, error) {
	return nil, nil
}

func (n *Node) Notify(ctx context.Context, empty *chord.Empty) (*chord.Alive, error) {
	return nil, nil
}

func (n *Node) Shutdown(ctx context.Context, empty *chord.Empty) (*chord.Empty, error) {
	k.node.Shutdown()
}


func

