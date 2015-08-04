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
	log "github.com/Sirupsen/logrus"
	"github.com/jordansissel/fingerpoken/zws"
	czmq "github.com/jordansissel/goczmq"
	"net/http"
)

func main() {
	endpoint := "inproc://fancy-req"
	go RunRep(endpoint)
	RunHTTP(":8111", "/zws/1.0")
}

func RunRep(endpoint string) {
	rep, err := czmq.NewRep(endpoint)
	if err != nil {
		panic(err)
	}

	for {
		message, flags, err := rep.RecvFrame()
		if err != nil {
			log.Fatalf("rep.RecvMessage failed: %s", err)
		}

		// Echo it back.
		err = rep.SendFrame(message, flags)
		if err != nil {
			log.Fatalf("rep.SendMessage failed: %s", err)
		}
	}
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
