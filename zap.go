package main

import (
	"fmt"
	czmq "github.com/zeromq/goczmq"
	"log"
)

func auth(s *czmq.Sock) {
	for {
		log.Printf("Auth: Waiting to receive\n")
		m, err := s.RecvMessage()
		if err != nil {
			log.Printf("AUTH RecvMessage: %s\n", err)
			return
		}
		for i, x := range m {
			log.Printf("AUTH auth %d: %s | %v\n", i, string(x), x)
		}

	}
}

func main() {
	var err error

	authsock := czmq.NewSock(czmq.Rep)
	authsock.Bind("inproc://zeromq.zap.01")
	go auth(authsock)

	//a := czmq.NewAuth()
	//a.Plain("/tmp/p/foo")

	server := czmq.NewSock(czmq.Rep)
	server.SetZapDomain("global")
	server.SetPlainServer(1)
	port, err := server.Bind("tcp://*:*")
	if err != nil {
		log.Fatalf("server.Bind failed: %s\n", err)
	}
	server_endpoint := fmt.Sprintf("tcp://127.0.0.1:%d", port)
	log.Printf("Server bound to %s", server_endpoint)

	client := czmq.NewSock(czmq.Req)
	client.SetPlainUsername("hello")
	client.SetPlainPassword("word")
	err = client.Connect(server_endpoint)
	if err != nil {
		log.Fatalf("client.SendFrame failed: %s\n", err)
	}
	log.Printf("Client connected to %s", server_endpoint)

	message := []byte("hello world")

	err = client.SendFrame([]byte("hello world"), 0)
	if err != nil {
		log.Fatalf("client.SendFrame failed: %s\n", err)
	}

	frame, _, err := server.RecvFrame()
	if err != nil {
		log.Fatalf("server.RecvFrame failed: %s\n", err)
	}

	err = server.SendFrame(frame, 0)
	if err != nil {
		log.Fatalf("server.SendFrame failed: %s\n", err)
	}

	frame, _, err = client.RecvFrame()
	if err != nil {
		log.Fatalf("client.RecvFrame failed: %s\n", err)
	}

	if string(message) != string(frame) {
		log.Printf("Messages did not match")
	} else {
		log.Fatalf("Messages matched! :)")
	}

}
