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
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/jordansissel/fingerpoken/mdp"
	"github.com/jordansissel/fingerpoken/zap"
	"github.com/jordansissel/fingerpoken/zws"
	czmq "github.com/zeromq/goczmq"
	"net/http"
)

func main() {
	agent, _ := zap.NewZapAgent()
	go agent.Run(&zap.OpenAccess{})
	defer agent.Destroy()

	broker_endpoint := "inproc://fancy-req"
	service := "echo"

	server_cert := czmq.NewCert()

	b, err := mdp.NewBroker(broker_endpoint)
	if err != nil {
		log.Fatalf("NewBroker(%s) failed: %s", broker_endpoint, err)
	}
	b.CurveCertificate = server_cert
	b.Bind("inproc://fancy-dealer")
	port, err := b.Bind("tcp://*:*")
	log.Printf("Broker available at binding: tcp://127.0.0.1:%d", port)
	go b.Run()

	// Start a worker.
	local_endpoint := fmt.Sprintf("tcp://127.0.0.1:%d", port)
	w := mdp.NewWorker(local_endpoint, service)
	w.CurveServerPublicKey = server_cert.PublicText()
	go w.Run(&echoWorker{})

	RunHTTP(":8111", "/zws/1.0")
}

func RunWorker(endpoint string) {
}

func RunHTTP(address, path string) {
	mux := http.NewServeMux()
	mux.HandleFunc(path, zws.HandleZWS)
	log.Printf("Listening on http://%s%s", address, path)
	server := http.Server{Addr: address, Handler: mux}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("server.ListenAndServe failed: %s", err)
	}
}
