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
  "github.com/jordansissel/fingerpoken/zws"
  "log"
  "net/http"
  "github.com/jordansissel/fingerpoken/mdp"
)

func main() {
  broker_endpoint := "inproc://fancy-req"
  service := "echo"

  b, err := mdp.NewBroker(broker_endpoint)
  if err != nil {
    log.Fatalf("NewBroker(%s) failed: %s", broker_endpoint, err)
  }
  port, err := b.Bind("tcp://*:*")
  log.Printf("Broker available at binding: tcp://x:%d", port)
  go b.Run()

  // Start a worker.
	w := mdp.NewWorker(broker_endpoint, service)
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
