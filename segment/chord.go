package main

import (
	"crypto/sha1"
	"io"
	"log"
	"net"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/context"

	"github.com/hoffa2/worm/protobuf/chord"
)

type Node interface {
	Shutdown()
	IsReachable(string) bool
	GetInactiveHost(map[string]chord.Node) (string, error)
	StartSegment(string) error
	GetAllReachable() []string
}

type lease struct {
	chord.Node
	time.Time
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
	logError          *log.Logger
	logInfo           *log.Logger
	logfile           *os.File
	leaseLock         sync.Mutex
	leases            leaseSlice
	aliveChan         chan chord.Node
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

	f, err := os.OpenFile("/var/log/jooonnaslog", os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	infolog := log.New(f, "\x1b[32m"+addr+"\x1b[0m"+" --> ", log.Lshortfile)
	errlog := log.New(f, "\x1b[31m"+addr+"\x1b[0m"+" --> ", log.Lshortfile)

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
		NetworkSizeTarget: 1,
		aliveSegments:     make(map[string]chord.Node),
		logError:          errlog,
		logInfo:           infolog,
		logfile:           f,
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
				h.logError.Println(err)
				h.electNewLeader()
				return
			}

			new := &chord.NewTarget{
				From: oldTarget,
				To:   newTarget,
			}

			alive, err := conn.Notify(ctx, new)
			if err != nil {
				h.logError.Println(err)
				h.electNewLeader()
				return
			}
			atomic.StoreInt32(&h.NetworkSizeTarget, alive.GetTarget())
		}()
	}
	if diff > 0 {
		go h.startupNodes(int(diff))
	} else if diff < 0 {
		go h.shutdownNodes(int(diff))
	}
}

// Alive is only issued to the leader
func (h *Heartbeat) Alive(ctx context.Context, alive *chord.Alive) (*chord.Alive, error) {

	h.aliveChan <- *alive.Node

	return &chord.Alive{IsAlive: true, Target: atomic.LoadInt32(&h.NetworkSizeTarget)}, nil
}

func (h *Heartbeat) Notify(ctx context.Context, new *chord.NewTarget) (*chord.Alive, error) {
	h.state.Lock()
	defer h.state.Unlock()

	if h.isLeader {
		h.UpdateNetworkSizeTarget(new.To)
	} else {

	}

	return &chord.Alive{Target: atomic.LoadInt32(&h.NetworkSizeTarget)}, nil
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
			h.logError.Println(err)
			return
		}
		h.upperLayer.StartSegment(host)

		h.aliveSegments[host] = h.CreateNodeAttr(host)
	}
}

func (h *Heartbeat) AbandonShip() {

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
			h.electNewLeader()
			break
		}
	}
}

func (h *Heartbeat) LeaderHint(leader string) {
	h.leader = h.CreateNodeAttr(leader)
	alive := h.NotifyAlive()
	if !alive {
		go h.ProbeNetwork()
	}
}

func (h *Heartbeat) updateLease(node *chord.Node) {
	h.leaseLock.Lock()
	defer h.leaseLock.Unlock()

	lease := lease{Node: *node, Time: time.Now()}

	h.leases = append(h.leases, lease)
	sort.Sort(h.leases)
}

func (h *Heartbeat) adhereLeases() {
	for {
		aliveSegments := make(map[string]chord.Node)
		for {
			select {
			case aliveSegment := <-h.aliveChan:
				aliveSegments[aliveSegment.IpAddress] = aliveSegment
			case <-time.After(time.Millisecond * time.Duration((len(h.aliveSegments) * 10))):
				break
			}
		}

		oldFucker := atomic.LoadInt32(&h.NetworkSizeTarget)
		diff := len(aliveSegments) - int(oldFucker)
		h.aliveSegments = aliveSegments
		if diff != 0 {
			if diff < 0 {
				go h.startupNodes(diff)
			} else {
				go h.shutdownNodes(diff)
			}
		}
	}
}

func (h *Heartbeat) handler() {
	for {
		// If we are not the leader: update lease
		// If we are the leader: collect leases
	}
}

func (h *Heartbeat) electNewLeader() {
	var segments segmentSlice

	for _, node := range h.aliveSegments {
		segments = append(segments, node)
	}

	sort.Sort(segments)

	// I'm the new leader
	if segments[0].ID == h.self.ID {
		h.leader = *h.self
		return
	}

	// Find new leader by pinging the nodes
	// in descending order. The first one
	// that responds is the new leader
	for i := 0; i < len(segments); i++ {
		h.leader = segments[i]
		alive := h.NotifyAlive()
		if alive {
			break
		}
	}
}

// Sort interface for leases
type leaseSlice []lease

func (l leaseSlice) Len() int {
	return len(l)
}

func (l leaseSlice) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l leaseSlice) Less(i, j int) bool {
	return time.Since(l[i].Time) < time.Since(l[j].Time)
}

// Sort interface for nodes
type segmentSlice []chord.Node

func (l segmentSlice) Len() int {
	return len(l)
}

func (l segmentSlice) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l segmentSlice) Less(i, j int) bool {
	return l[i].ID < l[j].ID
}
