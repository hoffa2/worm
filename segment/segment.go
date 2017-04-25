package segment

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
	"sync/atomic"
	"syscall"

	"github.com/urfave/cli"
)

var wormgatePort string
var segmentPort string

var hostname string

var targetSegments int32

type segment struct {
	// Underlying listener
	// for the http server
	net.Listener
	*Kademlia
	*RpcServer
}

func Run(c *cli.Context) error {
	host := c.String("host")
	segmentPort := c.String("segmentport")
	wormgatePort := c.String("wormgateport")
	if mode == "spread" {
		return SendSegment(host, segmentPort, wormgatPort)
	} else if mode == "start" {
		return StartSegmentServer(c)
	}
	return errors.New("Mode must be either spread or start")
}

func SendSegment(host, segmentPort, wormgatePort string) error {

	url := fmt.Sprintf("http://%s%s/wormgate?sp=%s", host, wormgatePort, segmentPort)
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

func StartSegmentServer(c *cli.Context) error {
	segmentPort := c.String("segmentport")

	// Startup case
	// Only bootstrap segments if "target" is set

	srv := http.Server{}

	l, err := net.Listen("tcp", ":"+segmentPort)
	if err != nil {
		log.Fatal(err)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	segment := &segment{
		Listener: l,
	}

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
	log.Printf("Reachable hosts: %s", strings.Join(fetchReachableHosts(), " "))

	if c.IsSet("target") {
		bootstrapNodes(c.Int("target"))
	} else {

	}

	go segment.runSegmentUntilShutdown()

	err = srv.Serve(segment.Listener)
	if err != nil {
		return err
	}
	return nil
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
	atomic.StoreInt32(&targetSegments, ts)
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

	// Shut down
	log.Printf("Received shutdown command, committing suicide")
	os.Exit(0)
}

func fetchReachableHosts() []string {
	url := fmt.Sprintf("http://localhost%s/reachablehosts", wormgatePort)
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

}

func (s *segment) Shutdown() {
	s.RpcServer.GracefulStop()
	s.Listener.Close()
}

func (s *segment) bootstrapNodes(int target) {
	for i := 0; i < targert; i++ {
		sendSegment()
	}
}
