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
package zws

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	czmq "github.com/jordansissel/goczmq"
	"net/http"
	"sync"
	"testing"
)

func TestZWS(t *testing.T) {
	wg := sync.WaitGroup{}
	// Backend
	rep, _ := czmq.NewRep("inproc://fancy-req")
	wg.Add(1)
	go func() {
		wg.Done()
		for {
			log.Printf("Rep socket ready to receive")
			message, flags, err := rep.RecvFrame()
			if err != nil {
				log.Fatalf("rep.RecvMessage failed: %s", err)
			}
			log.Printf("rep received: %s\n", string(message))
			// Echo it back.
			err = rep.SendFrame(message, flags)
			if err != nil {
				log.Fatalf("rep.SendMessage failed: %s", err)
			}
		}
	}()

	// ZWS HTTP
	mux := http.NewServeMux()
	mux.HandleFunc("/zws/1.0", HandleZWS)
	server := http.Server{Addr: ":12345", Handler: mux}
	wg.Add(1)
	go func() {
		wg.Done()
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("server.ListenAndServe failed: %s", err)
		}
	}()

	wg.Wait()

	// Client
	dialer := websocket.Dialer{Subprotocols: []string{"ZWS1.0"}}
	conn, _, err := dialer.Dial("ws://localhost:12345/zws/1.0?type=req", http.Header{})
	if err != nil {
		t.Errorf("websocket dial failed: %s", err)
		return
	}

	payload := []byte{'0'}
	greeting := "Hello world"
	payload = append(payload, []byte(greeting)...)
	err = conn.WriteMessage(websocket.TextMessage, payload)
	if err != nil {
		t.Errorf("conn.WriteMessage failed; %s", err)
		return
	}

	messageType, payload, err := conn.ReadMessage()
	if err != nil {
		t.Errorf("conn.ReadMessage failed; %s", err)
		return
	}

	if messageType != websocket.TextMessage {
		t.Errorf("conn.ReadMessage mesage type was not TextMessage")
		return
	}

	log.Printf("Payload; %s\n", string(payload))
}
