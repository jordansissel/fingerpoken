// This file is part of fingerpoken
// Copyright (C) 2015 Jordan Sissel
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
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

type Notification struct {
	Message string `json:"message"`
	Subtext string `json:"subtext"`
	Icon    string `json:"icon"`
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

	defer func() {
		for _, e := range endpoints {
			e.DestroyClose()
		}
	}()

	stdin := bufio.NewReader(os.Stdin)
	for {
		line := Notification{Icon: "info"}
		payload, err := stdin.ReadBytes('\n')
		if len(payload) == 0 {
			break
		}
		line.Message = string(payload[0 : len(payload)-1])
		log.Printf("Sending: %s\n", string(line.Message))
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
