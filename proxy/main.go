package main

import (
	"fmt"
	"github.com/gorilla/mux"
	consul "github.com/hashicorp/consul/api"
	"github.com/jordansissel/fingerpoken/util"
	czmq "github.com/zeromq/goczmq"
	"os"
	//"net"
	"log"
	"net/http"
)

type Notification []byte

type Gateway struct{}

func (g *Gateway) Ping(args interface{}, reply *interface{}) (err error) {
	return nil
}

func main() {
	log.SetOutput(os.Stdout)
	client, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		log.Fatalf("Failure to get a consul client connection: %s\n", err)
	}

	notification_chan := make(chan *Notification)
	go RunWebInterface(client, notification_chan)
	go RunRPCInterface(client)
	RunNotificationReceiver(client, notification_chan)
}

func RunWebInterface(client *consul.Client, notification_chan chan *Notification) {
	r := mux.NewRouter()
	r.Handle("/", http.FileServer(http.Dir("./static")))
	r.PathPrefix("/js/").Handler(http.StripPrefix("/js", http.FileServer(http.Dir("./static/js/"))))
	r.HandleFunc("/ws",
		func(w http.ResponseWriter, r *http.Request) {
			WebSocketMuxHandler(w, r, notification_chan)
		})
	http.Handle("/", r)
	http.ListenAndServe(":8000", nil)
}

func RunRPCInterface(client *consul.Client) {
	zj, err := util.NewZJServer()
	if err != nil {
		fmt.Printf("NewZJServer failure: %s\n", err)
		panic("!")
	}
	zj.RegisterWithConsul(client)
	zj.Register(&Gateway{})
	err = zj.Loop()
	fmt.Printf("Loop: %s\n", err)
}

func RunNotificationReceiver(client *consul.Client, notification_chan chan *Notification) {
	socket, err := czmq.NewPull("tcp://*:*")
	if err != nil {
		log.Printf("czmq.NewPull() failed: %s\n", err)
		panic("czmq.NewPull")
	}
	endpoint := socket.LastEndpoint()
	log.Printf("Notifications endpoint: %s\n", endpoint)
	err = util.ConsulRegisterService(client, "gateway", endpoint)
	log.Printf("ConsulRegisterService: %s\n", err)

	for {
		message, err := socket.RecvMessage()
		if err != nil {
			log.Printf("PULL: socket.RecvMessage(): %s\n", err)
			continue
		}
		log.Printf("PULL: Received: %s\n", string(message[0]))
		log.Printf("PULL: CHAN %v\n", notification_chan)
		var n = Notification(message[0])
		notification_chan <- &n
	}
}
