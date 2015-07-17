package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	consul "github.com/hashicorp/consul/api"
	czmq "github.com/zeromq/goczmq"
	"log"
	"os"
)

type Line struct {
	Payload string `json:"line"`
}

type JSONRPCNotification struct {
	Id     interface{} `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

func main() {
	log.SetOutput(os.Stdout)
	client, _ := consul.NewClient(consul.DefaultConfig())
	services, _, _ := client.Catalog().Service("gateway", "", nil)

	var endpoints []*czmq.Sock
	for _, service := range services {
		endpoint := fmt.Sprintf("tcp://%s:%d", service.Address, service.ServicePort)
		log.Println(endpoint)
		push, err := czmq.NewPush(endpoint)
		if err != nil {
			log.Printf("Failure setting up new PUSH to %s: %s\n", endpoint, err)
			continue
		}
		endpoints = append(endpoints, push)
	}

	stdin := bufio.NewReader(os.Stdin)
	for {
		line := Line{}
		payload, err := stdin.ReadBytes('\n')
		line.Payload = string(payload)
		log.Printf("Sending: %s\n", string(line.Payload))
		if err != nil {
			log.Printf("Error: %s\n", err)
			break
		}

		n := JSONRPCNotification{Id: nil, Method: "text", Params: line}
		blob, err := json.Marshal(n)
		if err != nil {
			log.Printf("Error serializing to json: %s\n", err)
			continue
		}

		log.Printf("%s\n", string(blob))
		for _, endpoint := range endpoints {
			endpoint.SendFrame(blob, 0)
		}
	}
}
