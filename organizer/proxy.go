package organize

import "net"

type Proxy struct {
	in *net.Conn
}

func SetupProxy(port string) (*Proxy, error) {
	conn, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		return nil, err
	}

	proxy := &Proxy{in: conn}

	go proxy.handleTraffic()
}

func (p *Proxy) handleTraffic() {
	for {
		conn.
	}
}
