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
	"github.com/jordansissel/fingerpoken/mdp"
	"github.com/jordansissel/fingerpoken/zap"
	"github.com/jordansissel/fingerpoken/zws"
	czmq "github.com/zeromq/goczmq"
	"log"
	"net/http"
	"time"
)

//go func() {
//client := mdp.NewClient(broker_endpoint)
//for {
//start := time.Now()
//response, err := client.SendRecv("webclient", [][]byte{[]byte(`{ "method": "openurl", "params": [ { "url": "https://elastic.zoom.us/j/2781360617" } ], "id": "1234" }`)})
//log.Printf("Response: %s", response)
//log.Printf("Err: %s", err)
//log.Printf("Duration: %s", time.Since(start))
//time.Sleep(1 * time.Second)
//}
//}()
