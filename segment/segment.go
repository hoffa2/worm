package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/hoffa2/worm/protobuf/chord"
	"github.com/soheilhy/cmux"
	"github.com/urfave/cli"
)

var (
	ErrNoReachableHosts = errors.New("No reachableHosts")
)

type segment struct {
	// Underlying listener
	// for the http server
	net.Listener
	heartbeat      *Heartbeat
	hostsLock      sync.RWMutex
	reachableHosts []string
	wormgatePort   string
	segmentPort    string
	hostname       string
}

func main() {
	app := cli.NewApp()
	app.Name = "Wormgate"
	app.Usage = "Run one of the components"
	app.Commands = []cli.Command{
		{
			Name:  "segment",
			Usage: "run segment",
			Action: func(c *cli.Context) error {
				if !c.IsSet("mode") {
					return errors.New("Wormport flag must be set")
				}
				return Run(c)
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "wormport, wp",
					Usage: "Wormagte port (prefix with colon)",
				},
				cli.StringFlag{
					Name:  "segmentport, sp",
					Usage: "segment port (prefix with colon)",
				},
				cli.StringFlag{
					Name:  "mode, m",
					Usage: "Spread or Start",
				},
				cli.StringFlag{
					Name:  "remotehost, rh",
					Usage: "Spread or Start",
				},
				cli.IntFlag{
					Name:  "target, t",
					Usage: "Inital number of targets (Set only if the segments is the first in the network)",
				},
			},
		},
	}
}

func Run(c *cli.Context) error {
	host := c.String("host")
	segmentPort := c.String("segmentport")
	wormgatePort := c.String("wormgateport")
	mode := c.String("mode")
	if mode == "spread" {
		return SendSegment(host, segmentPort, wormgatePort, "")
	} else if mode == "start" {
		return StartSegmentServer(c)
	}
	return errors.New("Mode must be either spread or start")
}

func SendSegment(host, segmentPort, wormgatePort, remoteHost string) error {

	url := fmt.Sprintf("http://%s%s/wormgate?sp=%s&rh=%s", host, wormgatePort, segmentPort, remoteHost)
	filename := "tmp.tar.gz"

	log.Printf("Spreading to %s", url)

	// ship the binary and the qml file that describes our screen output
	tarCmd := exec.Command("tar", "-zc", "-f", filename, "segment")
	tarCmd.Run()
	defer os.Remove(filename)

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Could not read input file %s", err.Error())
	}

	resp, err := http.Post(url, "string", file)
	if err != nil {
		return fmt.Errorf("POST error %s", err.Error())
	}

	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()

	if resp.StatusCode == 200 {
		log.Println("Received OK from server")
	} else {
		log.Println("Response: ", resp)
	}
	return nil
}

func (s *segment) StartSegment(host string) error {
	return SendSegment(host, s.segmentPort, s.wormgatePort, s.hostname)
}

func StartSegmentServer(c *cli.Context) error {
	segmentPort := c.String("segmentport")
	target := c.Int("target")
	wormgatePort := c.String("wormgateport")
	remoteHost := c.String("remotehost")

	// Startup case
	// Only bootstrap segments if "target" is set

	srv := http.Server{}

	l, err := net.Listen("udp", ":"+segmentPort)
	if err != nil {
		log.Fatal(err)
	}

	hostname, _ := os.Hostname()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	m := cmux.New(l)

	grpcL := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpL := m.Match(cmux.HTTP1Fast())

	segment := &segment{
		Listener:     l,
		segmentPort:  segmentPort,
		wormgatePort: wormgatePort,
		hostname:     hostname,
	}
	segment.heartbeat, err = SetupHeartbeat(hostname, segmentPort, segment, grpcL)
	if err != nil {
		return err
	}

	segment.heartbeat.UpdateNetworkSizeTarget(int32(target))
	segment.heartbeat.LeaderHint(remoteHost)

	// Making sure that the port
	// is closed if we are killed
	go func() {
		_ = <-sigc
		segment.Shutdown()
	}()

	http.HandleFunc("/", segment.IndexHandler)
	http.HandleFunc("/targetsegments", segment.targetSegmentsHandler)
	http.HandleFunc("/shutdown", segment.shutdownHandler)

	log.Printf("Starting segment server on %s%s\n", hostname, segmentPort)
	//log.Printf("Reachable hosts: %s", strings.Join(fetchReachableHosts(), " "))

	go srv.Serve(httpL)
	return m.Serve()
}

// IndexHandler handles lol
func (s segment) IndexHandler(w http.ResponseWriter, r *http.Request) {

	// We don't use the request body. But we should consume it anyway.
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()

	killRateGuess := 2.0

	fmt.Fprintf(w, "%.3f\n", killRateGuess)
}

func (s segment) targetSegmentsHandler(w http.ResponseWriter, r *http.Request) {

	var ts int32
	pc, rateErr := fmt.Fscanf(r.Body, "%d", &ts)
	if pc != 1 || rateErr != nil {
		log.Printf("Error parsing targetSegments (%d items): %s", pc, rateErr)
	}

	// Consume and close rest of body
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()

	log.Printf("New targetSegments: %d", ts)
	s.heartbeat.UpdateNetworkSizeTarget(ts)
}

func (s segment) shutdownHandler(w http.ResponseWriter, r *http.Request) {

	// Consume and close body
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()

	// Closing the underlying listener
	// for the http server
	err := s.Listener.Close()
	if err != nil {
		log.Println(err)
	}

	s.heartbeat.AbandonShip()

	// Shut down
	log.Printf("Received shutdown command, committing suicide")
	os.Exit(0)
}

func (s segment) fetchReachableHosts() []string {
	url := fmt.Sprintf("http://localhost%s/reachablehosts", s.wormgatePort)
	resp, err := http.Get(url)
	if err != nil {
		return []string{}
	}

	var bytes []byte
	bytes, err = ioutil.ReadAll(resp.Body)
	body := string(bytes)
	resp.Body.Close()

	trimmed := strings.TrimSpace(body)
	nodes := strings.Split(trimmed, "\n")
	return nodes
}

func (s *segment) runSegmentUntilShutdown() {
	for {
		hosts := s.fetchReachableHosts()
		s.hostsLock.Lock()
		s.reachableHosts = hosts
		s.hostsLock.Unlock()
		time.Sleep(time.Second * 1)
	}
}

func (s *segment) IsReachable(host string) bool {
	for _, h := range s.reachableHosts {
		if h == host {
			return true
		}
	}

	return false
}

func (s *segment) Shutdown() {
	s.Listener.Close()
}

func (s *segment) GetInactiveHost(active map[string]chord.Node) (string, error) {
	s.hostsLock.RLock()
	defer s.hostsLock.RUnlock()

	for _, host := range s.reachableHosts {
		if _, ok := active[host]; ok {
			return host, nil
		}
	}
	return "", ErrNoReachableHosts
}

func (s *segment) GetAllReachable() []string {
	return s.reachableHosts
}
