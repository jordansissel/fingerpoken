package main

import (
	"fmt"
	"github.com/gorilla/mux"
	consul "github.com/hashicorp/consul/api"
	"github.com/jordansissel/fingerpoken/util"
	czmq "github.com/jordansissel/goczmq"
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

  nbc := &NotificationBroadcaster{}
	go RunWebInterface(client, nbc)
	go RunRPCInterface(client)
	RunNotificationReceiver(client, nbc)
}

func RunWebInterface(client *consul.Client, nbc *NotificationBroadcaster) {
	r := mux.NewRouter()
	r.Handle("/", http.FileServer(http.Dir("./static")))
	r.PathPrefix("/js/").Handler(http.StripPrefix("/js", http.FileServer(http.Dir("./static/js/"))))
	r.HandleFunc("/ws",
		func(w http.ResponseWriter, r *http.Request) {
			WebSocketMuxHandler(w, r, nbc)
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

func RunNotificationReceiver(client *consul.Client, nbc *NotificationBroadcaster) {
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
		//log.Printf("PULL: Received: %s\n", string(message[0]))
		var n = Notification(message[0])
    nbc.Publish(&n);
		//log.Printf("PULL: CHAN %v\n", nbc)
	}
}
