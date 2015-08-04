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
	flags "github.com/jessevdk/go-flags"
	"github.com/jordansissel/fingerpoken/mdp"
	"github.com/jordansissel/fingerpoken/zap"
	"github.com/jordansissel/fingerpoken/zws"
	czmq "github.com/zeromq/goczmq"
	"net/http"
	"os"
)

type Settings struct {
	ServerCertificatePath string `long:"server-certificate" required:"true"`
	// TODO(sissel): Allow specifyinj the port to bind to.
}

func main() {
	var settings Settings
	_, err := flags.Parse(&settings)
	if err != nil {
		switch err.(*flags.Error).Type {
		case flags.ErrHelp:
			os.Exit(0)
		default:
			log.Error(err.Error())
			os.Exit(1)
		}
	}

	agent, _ := zap.NewZapAgent()
	// TODO(sissel): Gate access to prevent untrusted connections.
	go agent.Run(&zap.OpenAccess{})
	defer agent.Destroy()

	broker_endpoint := "inproc://fancy-req"
	server_cert, err := czmq.NewCertFromFile(settings.ServerCertificatePath)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	b, err := mdp.NewBroker(broker_endpoint)
	if err != nil {
		log.Fatalf("NewBroker(%s) failed: %s", broker_endpoint, err)
	}
	b.CurveCertificate = server_cert
	b.Bind("inproc://fancy-dealer")
	port, err := b.Bind("tcp://*:*")
	log.WithFields(log.Fields{"address": fmt.Sprintf("tcp://*:%d", port)}).Info("Broker available")
	go RunHTTP(":8111", "/zws/1.0")
	b.Run()
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
