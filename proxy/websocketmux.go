package main

import (
	"github.com/gorilla/websocket"
	//consul "github.com/hashicorp/consul/api"
	//czmq "github.com/zeromq/goczmq"
	"log"
	"net"
	"net/http"
	//"fmt"
)

type WebSocketMux struct {
}

// The return value of websocket.ReadMessage wrapped in a struct.
type WebSocketMessage struct {
	messageType int
	p           []byte
}

func WebSocketMuxHandler(w http.ResponseWriter, r *http.Request, notifications chan *Notification) {
	mux := &WebSocketMux{}

	mux.Handle(w, r, notifications)
}

func upgrade(w http.ResponseWriter, r *http.Request) (conn *websocket.Conn, err error) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (wsm *WebSocketMux) Handle(w http.ResponseWriter, r *http.Request, notifications chan *Notification) {
	conn, err := upgrade(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("New websocket connection\n")

	var tcp = conn.UnderlyingConn().(*net.TCPConn)
	tcp.SetNoDelay(true)

	var input chan *WebSocketMessage

	// Notifications come in on the notifications chan. Forward them.
	go WebSocketForwardNotifications(conn, notifications)
	go WebSocketReadMessagesLoop(conn, input)

	for x := range input {
		log.Printf("ReadMessage: %s\n", string(x.p))
	}
}

func WebSocketReadMessagesLoop(conn *websocket.Conn, input chan *WebSocketMessage) {
	for {
		mtype, payload, err := conn.ReadMessage()
		if err != nil {
			log.Printf("conn.ReadMessages err: %s\n", err)
			break
		}
		log.Printf("WebSocket Message: %s\n", string(payload))
		input <- &WebSocketMessage{mtype, payload}
	}
	conn.Close()
}

func WebSocketForwardNotifications(conn *websocket.Conn, notifications chan *Notification) {
	for {
		notification := <-notifications
		log.Printf("Received %v\n", string(*notification))
		err := conn.WriteMessage(1, *notification)
		if err != nil {
			log.Printf("Error writing to websocket %v\n", *conn)
			break
		}
	}
}

/*
func x() {

	client, _ := consul.NewClient(consul.DefaultConfig())
	services, _, _ := client.Catalog().Service("rpc", "", nil)

	var endpoints []*czmq.Sock
	for _, service := range services {
		endpoint := fmt.Sprintf("tcp://%s:%d", service.Address, service.ServicePort)
		log.Println(endpoint)
		req, _ := czmq.NewReq(endpoint)
		endpoints = append(endpoints, req)
	}

  err := nil
	for err == nil {
		mtype, payload, err := conn.ReadMessage()
		if err != nil {
			break
		}
		//log.Printf("%s Payload [%v]: %s\n", mtype, err, string(payload))
		for _, endpoint := range endpoints {
			err = endpoint.SendFrame(payload[:], 0)
			if err != nil {
				log.Printf("endpoint.SendMessage fail: %s\n", err)
				panic("!")
			}
			response, err := endpoint.RecvMessage()
			if err != nil {
				log.Printf("endpoint.RecvMessage fail: %s\n", err)
				panic("!")
			}

			//time.Sleep(50 * time.Millisecond)
			//log.Printf("Resp: %v\n", string(response[0]))
			err = conn.WriteMessage(mtype, response[0])
			if err != nil {
				conn.Close()
			}
		}
	}
}

*/