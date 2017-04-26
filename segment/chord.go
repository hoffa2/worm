package segment

import (
	"crypto/sha1"
	"golang.org/x/net/context"
	"io"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hoffa2/worm/protobuf/chord"
)

type Node interface {
	Shutdown()
	IsReachable(string) bool
	GetInactiveHost(map[string]chord.Node) (string, error)
	StartSegment(string) error
	GetAllReachable() []string
}

type Heartbeat struct {
	self       *chord.Node
	upperLayer Node
	*RpcServer
	*ClientRemote
	NetworkSizeTarget int32
	leader            chord.Node
	aliveSegments     map[string]chord.Node
	isLeader          bool
	startupLock       sync.Mutex
	state             sync.Mutex
}

func hashValue(str string) string {
	h := sha1.New()
	io.WriteString(h, str)
	return string(h.Sum(nil))
}

func SetupHeartbeat(addr, port string, nodeinterface Node, l net.Listener) (*Heartbeat, error) {
	h, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	n := &chord.Node{
		ID:        hashValue(h),
		IpAddress: addr,
		RpcPort:   port,
	}

	remote := SetupRemote(nodeinterface)

	heartbeat := &Heartbeat{
		self:              n,
		upperLayer:        nodeinterface,
		ClientRemote:      remote,
		NetworkSizeTarget: 0,
		aliveSegments:     make(map[string]chord.Node),
	}
	rpcServer, err := SetupRPCServer(heartbeat, port, l)
	if err != nil {
		return nil, err
	}

	heartbeat.RpcServer = rpcServer

	return heartbeat, nil
}

func (h *Heartbeat) UpdateNetworkSizeTarget(newTarget int32) {
	oldTarget := atomic.LoadInt32(&h.NetworkSizeTarget)
	diff := oldTarget - newTarget
	atomic.StoreInt32(&h.NetworkSizeTarget, newTarget)
	if !h.isLeader {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
			defer cancel()
			conn, err := h.GetConn(&h.leader)
			if err != nil {
				//Leader is dead
			}

			new := &chord.NewTarget{
				From: oldTarget,
				To:   newTarget,
			}

			alive, err := conn.Notify(ctx, new)
			if err != nil {
				//Leader is dead
			}
			h.UpdateNetworkSizeTarget(alive.GetTarget())
		}()
	}
	if diff < 0 {
		go h.startupNodes(int(diff))
	} else if diff > 0 {
		go h.shutdownNodes(int(diff))
	}
}

func (h *Heartbeat) Alive(ctx context.Context, alive *chord.Alive) (*chord.Alive, error) {
	return &chord.Alive{IsAlive: true}, nil
}

func (h *Heartbeat) Notify(ctx context.Context, new *chord.NewTarget) (*chord.Alive, error) {
	h.state.Lock()
	defer h.state.Unlock()

	h.UpdateNetworkSizeTarget(new.To)

	return nil, nil
}

func (h *Heartbeat) Shutdown(ctx context.Context, empty *chord.Empty) (*chord.Empty, error) {
	h.upperLayer.Shutdown()
	h.RpcServer.CloseServer()
	return &chord.Empty{}, nil
}

// Notify the leader that I'm not dead
func (h *Heartbeat) NotifyAlive() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	target := atomic.LoadInt32(&h.NetworkSizeTarget)
	conn, err := h.GetConn(&h.leader)
	if err != nil {
		return false
	}
	alive, err := conn.Alive(ctx, &chord.Alive{Target: target})
	if err != nil {
		return false
	}

	h.UpdateNetworkSizeTarget(alive.GetTarget())

	return true
}

func (h *Heartbeat) startupNodes(n int) {
	for i := 0; i < n; i++ {
		host, err := h.upperLayer.GetInactiveHost(h.aliveSegments)
		if err == ErrNoReachableHosts {
			return
		}
		h.upperLayer.StartSegment(host)
	}
}

func (h *Heartbeat) shutdownNodes(n int) {
	h.startupLock.Lock()
	defer h.startupLock.Unlock()
	killed := 0
	for key, segment := range h.aliveSegments {
		conn, err := h.GetConn(&segment)
		if err != nil {
			delete(h.aliveSegments, key)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
		defer cancel()
		// Don't care whether call fails or not
		// either way node is dead
		_, _ = conn.Shutdown(ctx, &chord.Empty{})

		killed++
		if killed == n {
			return
		}
	}
}

func (h *Heartbeat) CreateNodeAttr(host string) chord.Node {
	return chord.Node{
		ID:        hashValue(host),
		RpcPort:   h.self.RpcPort,
		IpAddress: host,
	}
}

func (h *Heartbeat) ProbeNetwork() {
	hosts := h.upperLayer.GetAllReachable()
	var wg sync.WaitGroup
	aliveChan := make(chan chord.Node, len(hosts))

	for _, host := range hosts {
		wg.Add(1)
		go func(host string, aliveChan chan chord.Node) {
			defer wg.Done()
			node := h.CreateNodeAttr(host)

			_, err := h.GetConn(&node)
			if err != nil {
				return
			}
			aliveChan <- node
		}(host, aliveChan)
	}
	wg.Wait()

	for {
		select {
		case node := <-aliveChan:
			h.aliveSegments[node.IpAddress] = node
		default:
			return
		}
	}
}
