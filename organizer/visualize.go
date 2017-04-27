package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/hoffa2/worm/grapher"
	"github.com/hoffa2/worm/protobuf"
	"github.com/hoffa2/worm/rocks"
	"github.com/urfave/cli"
)

const minx, maxx = 1, 3
const miny, maxy = 0, 59
const colwidth = 20
const refreshRate = 100 * time.Millisecond
const pollRate = refreshRate / 2
const pollErrWait = 20 * time.Second

var wormgatePort string
var segmentPort string

type status struct {
	wormgate  bool
	segment   bool
	err       bool
	rateGuess float32
	rateErr   error
}

var statusMap struct {
	sync.RWMutex
	m map[string]status
}

var killRate int32
var targetSegments int32
var partitionScheme int32

// Use separate clients for wormgates vs segments
//
// There is something about making connections to the same host at different
// ports that confuses the connection caching and reuse. If we just use the
// default Client with the default Transfer, the number of open connections
// balloons during polling until we can't connect anymore. But using separate
// clients for each port (but multiple hosts) works fine.
//
var wormgateClient *http.Client
var segmentClient *http.Client

func createClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{},
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "Wormgate"
	app.Usage = "Run one of the components"
	app.Commands = []cli.Command{
		{
			Name:  "viz",
			Usage: "run visualizer",
			Action: func(c *cli.Context) error {
				if !c.IsSet("wormport") {
					return errors.New("Wormport flag must be set")
				}
				if !c.IsSet("segmentport") {
					return errors.New("segmentport flag must be set")
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
			},
		},
	}
}

func Run(c *cli.Context) error {
	wormgatePort = c.String("wormport")
	segmentPort = c.String("segmentport")

	nodes := rocks.ListNodes()

	statusMap.m = make(map[string]status)
	for _, node := range nodes {
		statusMap.m[node] = status{}
	}

	targetSegments = 5

	segmentClient = createClient()
	wormgateClient = createClient()

	// Catch interrupt and quit
	interrupt := make(chan os.Signal, 2)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-interrupt
		fmt.Print(ansi_clear_to_end)
		fmt.Println()
		log.Print("Shutting down")
		os.Exit(0)
	}()

	// Start poll routines
	for node := range statusMap.m {
		go pollNodeForever(node)
	}

	grapher.Init(inputHandler, "8001")

	// Start random node killer
	go killNodesForever()

	// Loop display forever
	for {
		//printNodeGrid()
		time.Sleep(refreshRate)
	}
	return nil
}

func pollNodeForever(node string) {
	log.Printf("Starting poll routine for %s", node)
	for {
		s := pollNode(node)
		statusMap.Lock()
		statusMap.m[node] = s
		statusMap.Unlock()
		if s.err {
			time.Sleep(pollErrWait)
		} else {
			time.Sleep(pollRate)
		}
	}
}

func pollNode(host string) status {
	wormgateURL := fmt.Sprintf("http://%s%s/", host, wormgatePort)
	segmentURL := fmt.Sprintf("http://%s%s/", host, segmentPort)

	wormgate, _, wgerr := httpGetOk(wormgateClient, wormgateURL)
	if wgerr != nil {
		return status{false, false, true, 0, nil}
	}
	segment, segBody, segErr := httpGetOk(segmentClient, segmentURL)

	if segErr != nil {
		return status{false, false, true, 0, nil}
	}

	var rateGuess float32
	var rateErr error
	if segment {
		var pc int
		pc, rateErr = fmt.Sscanf(segBody, "%f", &rateGuess)
		if pc != 1 || rateErr != nil {
			log.Printf("Error parsing from %s (%d items): %s", host, pc, rateErr)
			log.Printf("Response %s: %s", host, segBody)
		}
	}

	return status{wormgate, segment, false, rateGuess, rateErr}
}

func httpGetOk(client *http.Client, url string) (bool, string, error) {
	resp, err := client.Get(url)
	isOk := err == nil && resp.StatusCode == 200
	body := ""
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "connection refused") {
			// ignore connection refused errors
			err = nil
		} else {
			log.Printf("Error checking %s: %s", url, err)
		}
	} else {
		var bytes []byte
		bytes, err = ioutil.ReadAll(resp.Body)
		body = string(bytes)
		resp.Body.Close()
	}
	return isOk, body, err
}

func changeSegmentTarget(oldTarget, newTarget int32) {
	if oldTarget != newTarget {
		for _, target := range randomSegment() {
			doTargetSegmentsPost(target, newTarget)
		}
	}
	atomic.StoreInt32(&targetSegments, newTarget)
}

func inputHandler(msg message.FromClient, send func(message.ToClient) error) error {
	switch msg.Msg.(type) {
	case *message.FromClient_ChangeTarget:
		newTs := msg.GetChangeTarget()
		ts := atomic.LoadInt32(&targetSegments)
		changeSegmentTarget(ts, newTs)
	case *message.FromClient_ShutdownTarget:
		shutdown := msg.GetShutdownTarget()
		if shutdown == true {
			for _, target := range randomSegment() {
				doWormShutdownPost(target)
			}
		}
	case *message.FromClient_GetTarget:
		if msg.GetGetTarget() == true {
			target := &message.ToClient{
				&message.ToClient_Target{
					Target: &message.Target{
						Target: atomic.LoadInt32(&targetSegments),
					},
				},
			}
			err := send(*target)
			if err != nil {
				fmt.Println(err)
			}
		}
	default:
		fmt.Println("LOL")
		break
	}

	return nil
}

func killNodesForever() {
	for {
		kr := atomic.LoadInt32(&killRate)
		if kr == 0 {
			// do nothing
			time.Sleep(time.Second)
		} else {
			killRandomNode()
			killWait := time.Duration(1000/kr) * time.Millisecond
			time.Sleep(killWait)
		}
	}
}

func randomSegment() []string {
	var segmentNodes []string
	statusMap.RLock()
	for node, status := range statusMap.m {
		if status.segment {
			segmentNodes = append(segmentNodes, node)
		}
	}
	statusMap.RUnlock()
	if len(segmentNodes) > 0 {
		ri := rand.Intn(len(segmentNodes))
		return segmentNodes[ri : ri+1]
	}
	return []string{}
}

func allWormgateNodes() []string {
	var nodes []string
	statusMap.RLock()
	for node, status := range statusMap.m {
		if status.wormgate {
			nodes = append(nodes, node)
		}
	}
	statusMap.RUnlock()
	return nodes
}

func killRandomNode() {
	for _, target := range randomSegment() {
		doKillPost(target)
	}
}

func doKillPost(node string) error {
	log.Printf("Killing segment on %s", node)
	url := fmt.Sprintf("http://%s%s/killsegment", node, wormgatePort)
	resp, err := wormgateClient.PostForm(url, nil)
	if err != nil && !strings.Contains(fmt.Sprint(err), "refused") {
		log.Printf("Error killing %s: %s", node, err)
	}
	if err == nil {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
	return err
}

func doPartitionSchemePost(node string, newps int32) error {
	log.Printf("Posting partitionScheme: %d -> %s", newps, node)

	url := fmt.Sprintf("http://%s%s/partitionscheme", node, wormgatePort)
	postBody := strings.NewReader(fmt.Sprint(newps))

	resp, err := wormgateClient.Post(url, "text/plain", postBody)
	if err != nil && !strings.Contains(fmt.Sprint(err), "refused") {
		log.Printf("Error posting partitionScheme %s: %s", node, err)
	}
	if err == nil {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
	return err
}

func doTargetSegmentsPost(node string, newts int32) error {
	log.Printf("Posting targetSegments: %d -> %s", newts, node)

	url := fmt.Sprintf("http://%s%s/targetsegments", node, segmentPort)
	postBody := strings.NewReader(fmt.Sprint(newts))

	resp, err := segmentClient.Post(url, "text/plain", postBody)
	if err != nil && !strings.Contains(fmt.Sprint(err), "refused") {
		log.Printf("Error posting shutdown to %s: %s", node, err)
	}
	if err == nil {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
	return err
}

func doWormShutdownPost(node string) error {
	log.Printf("Posting shutdown to %s", node)

	url := fmt.Sprintf("http://%s%s/shutdown", node, segmentPort)

	resp, err := segmentClient.PostForm(url, nil)
	if err != nil && !strings.Contains(fmt.Sprint(err), "refused") {
		log.Printf("Error posting targetSegments %s: %s", node, err)
	}
	if err == nil {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
	return err
}

const ansi_bold = "\033[1m"
const ansi_reset = "\033[0m"
const ansi_reverse = "\033[30;47m"
const ansi_red_bg = "\033[30;41m"
const ansi_clear_to_end = "\033[0J"

func ansi_down_lines(n int) string {
	return fmt.Sprintf("\033[%dE", n)
}
func ansi_up_lines(n int) string {
	return fmt.Sprintf("\033[%dF", n)
}

func printNodeGrid() {
	statusMap.RLock()

	gridBuf := bytes.NewBuffer(nil)
	rateGuesses := make([]float32, 0, len(statusMap.m))

	fmt.Fprint(gridBuf, ansi_clear_to_end)
	fmt.Fprintln(gridBuf)
	fmt.Fprintln(gridBuf)
	fmt.Fprint(gridBuf, "Legend: ")
	fmt.Fprint(gridBuf, "node,  ")
	fmt.Fprint(gridBuf, ansi_bold, "wormgate", ansi_reset, ",  ")
	fmt.Fprint(gridBuf, ansi_reverse, "segment", ansi_reset, ",  ")
	fmt.Fprint(gridBuf, ansi_red_bg, "error", ansi_reset)
	fmt.Fprintln(gridBuf)
	fmt.Fprint(gridBuf, "Keys  :")
	fmt.Fprint(gridBuf, "  kK/jJ kill rate,")
	fmt.Fprint(gridBuf, "  +/- segments,")
	fmt.Fprint(gridBuf, "  0-9 partition,")
	fmt.Fprint(gridBuf, "  s worm shutdown,")
	fmt.Fprint(gridBuf, "  Ctrl-C quit")

	for x := minx; x <= maxx; x++ {
		for y := miny; y <= maxy; y++ {
			if y%colwidth == 0 {
				fmt.Fprintf(gridBuf, "\n%d: %02d+", x, y/colwidth*colwidth)
			}
			if y%10 == 0 {
				fmt.Fprintf(gridBuf, "|")
			}
			node := fmt.Sprintf("compute-%d-%d", x, y)
			status, nodeup := statusMap.m[node]

			var char string
			if nodeup {
				char = fmt.Sprint(y % 10)
			} else {
				char = " "
			}

			if status.err {
				fmt.Fprint(gridBuf, ansi_red_bg)
			} else {
				if status.wormgate {
					fmt.Fprint(gridBuf, ansi_bold)
				}
				if status.segment {
					fmt.Fprint(gridBuf, ansi_reverse)
				}
				if status.segment && status.rateErr == nil {
					rateGuesses = append(rateGuesses,
						status.rateGuess)
				}
			}
			fmt.Fprint(gridBuf, char)
			fmt.Fprint(gridBuf, ansi_reset)
		}
	}
	statusMap.RUnlock()
	fmt.Fprintln(gridBuf)

	ts := atomic.LoadInt32(&targetSegments)
	fmt.Fprintf(gridBuf, "Target number of segments: %d\n", ts)

	kr := atomic.LoadInt32(&killRate)
	fmt.Fprintf(gridBuf, "Kill rate: %d/sec\n", kr)
	fmt.Fprintf(gridBuf, "Avg guess: %.1f/sec (%d segments reporting)\n",
		mean(rateGuesses), len(rateGuesses))

	fmt.Fprintln(gridBuf, time.Now().Format(time.StampMilli))
	var gridLines = bytes.Count(gridBuf.Bytes(), []byte("\n"))
	fmt.Fprint(gridBuf, ansi_up_lines(gridLines))
	io.Copy(os.Stdout, gridBuf)
}

func mean(floats []float32) float32 {
	var sum float32 = 0
	for _, f := range floats {
		sum += f
	}
	return sum / float32(len(floats))
}
