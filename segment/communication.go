package segment

import (
	"errors"
	"net"
)

var (
	ErrAddrInvalid = errors.New("Could not resolve addr")
)

// Gossip handling lol
type Gossip struct {
	activeNodes map[string]net.Conn
	net.Listener
}

// SetupGossip start listening for aliens
// param port port on which to accept connections
func SetupGossip(port string) (*Gossip, error) {
	conn, err := net.Listen("udp", "0.0.0.0:"+port)
	if err != nil {
		return nil, err
	}

	goss := &Gossip{
		Listener:    conn,
		activeNodes: make(map[string]net.Conn),
	}

	return goss, nil
}

// ProbeActive check for active nodes
// param segments nodes to probe
func (g Gossip) ProbeActive(segments []string) {
	for _, segment := range segments {
		if conn, ok := g.activeNodes[segment]; ok {
			_, err := conn.Write([]byte("LOL"))
			if err != nil {
				delete(g.activeNodes, segment)
			}

		}
	}
}
