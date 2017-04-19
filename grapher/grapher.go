package grapher

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/hoffa2/worm/protobuf"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Grapher struct {
	conn          *websocket.Conn
	onMessageFunc func(message.FromClient, func(message.ToClient) error) error
}

func (g Grapher) Addnode(nodeID int) error {
	msg := &message.Addnode{NodeId: int32(nodeID)}

	payload, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	err = g.conn.WriteMessage(websocket.BinaryMessage, payload)
	if err != nil {
		return err
	}

	return nil
}

func Init(onMessageFunc func(message.FromClient, func(message.ToClient) error) error, port string) *Grapher {
	grapher := &Grapher{
		onMessageFunc: onMessageFunc,
	}

	r := mux.NewRouter()

	gopath := os.Getenv("GOPATH")
	fmt.Println(gopath)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, gopath+"/src/github.com/hoffa2/worm/vis/index.html")
	})
	r.HandleFunc("/ws", grapher.onWSConn)
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets", http.FileServer(http.Dir(gopath+"/src/github.com/hoffa2/worm/vis"))))

	go func() {
		err := http.ListenAndServe("localhost:"+port, r)
		if err != nil {
			log.Println(err)
		}
	}()

	return grapher
}

func (g Grapher) onWSConn(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	g.conn = conn

	go g.handleTraffic()
}

func (g Grapher) handleTraffic() {
	var msg message.FromClient
	for {
		_, payload, err := g.conn.ReadMessage()
		if err != nil {
			return
		}
		log.Println("Got Message")

		err = proto.Unmarshal(payload, &msg)
		if err != nil {
			log.Println(err)
			return
		}

		g.onMessageFunc(msg, g.sendMessage)
		msg.Reset()
	}
}

func (g Grapher) sendMessage(msg message.ToClient) error {
	payload, err := proto.Marshal(&msg)
	if err != nil {
		return err
	}

	return g.conn.WriteMessage(websocket.BinaryMessage, payload)
}
