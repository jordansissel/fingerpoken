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
	"github.com/gorilla/websocket"
	czmq "github.com/zeromq/goczmq"
	"log"
	"net/http"
	"time"
)

const (
	FINAL_FRAME = '0'
	MORE_FRAME  = '1'
)

func upgrade(w http.ResponseWriter, r *http.Request) (conn *websocket.Conn, err error) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		Subprotocols:    []string{"ZWS1.0"},
		CheckOrigin: func(r *http.Request) bool {
			return true // TODO(sissel): Actually validate origin
		},
	}

	// TODO(sissel): This would be a good time to authenticate this request.
	// We could use Cookies or Authentication headers and ask ZAP if the request
	// should be accepted.

	conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func WStoZMQ(ws *websocket.Conn, zmq *czmq.Sock) {
	defer func() {
		deadline := time.Now().Add(1 * time.Second)
		ws.WriteControl(websocket.CloseMessage, []byte("Invalid MORE frame value (expected '0' or '1')"), deadline)
	}()

	for {
		messageType, message, err := ws.ReadMessage()

		log.Printf("received %d: %s\n", messageType, string(message))
		if err != nil {
			log.Printf("Error conn.ReadMessage: %s", err)
			ws.Close()
			return
		}
		if messageType != websocket.TextMessage {
			log.Printf("Got unexpected message type: %s", messageType)
			continue
		}

		flags := czmq.FlagNone
		switch more := message[0]; more {
		case FINAL_FRAME: // Nothing to do
		case MORE_FRAME: // There are more frames after this one.
			flags |= czmq.FlagMore
		default:
			log.Printf("Got an invalid MORE frame value: %c", more)
			break
		}

		err = zmq.SendFrame(message[1:], flags)
		if err != nil {
			log.Printf("zmq.SendFrame failed: %s\n", err)
			break
		}
	}
}

func ZMQtoWS(ws *websocket.Conn, zmq *czmq.Sock) {
	defer func() {
		deadline := time.Now().Add(1 * time.Second)
		ws.WriteControl(websocket.CloseMessage, []byte("Invalid MORE frame value (expected '0' or '1')"), deadline)
	}()

	for {
		frame, more, err := zmq.RecvFrame()
		if err != nil {
			log.Printf("zmq.RecvFrame failed: %s", err)
			break
		}

		log.Printf("z2w: received %s\n", string(frame))

		var payload [1]byte
		if more&czmq.FlagMore == czmq.FlagMore {
			payload[0] = MORE_FRAME
		} else {
			payload[0] = FINAL_FRAME
		}
		err = ws.WriteMessage(websocket.TextMessage, append(payload[:], frame...))
		if err != nil {
			log.Printf("ws.WriteMessage failed: %s", err)
			break
		}
	}
}

func HandleZWS(w http.ResponseWriter, r *http.Request) {
	endpoint_prefix := "inproc://fancy"
	zws, err := NewZWS(endpoint_prefix, w, r)
	//
	if err != nil {
		log.Printf("Problem creating new ZWS: %s", err)
		w.WriteHeader(http.StatusBadRequest)
	}

  defer zws.zmq.Destroy()

	switch zws.socket_type {
	case REQ:
		ProxyReqRep(zws.ws, zws.zmq)
	case SUB:
		ProxySub(zws.ws, zws.zmq)
	default:
		log.Printf("Invalid socket type, cannot handle.", err)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func ProxyReqRep(ws *websocket.Conn, zmq *czmq.Sock) {
	log.Printf("Proxy reqrep")
	defer func() {
		deadline := time.Now().Add(1 * time.Second)
		ws.WriteControl(websocket.CloseMessage, []byte("Invalid MORE frame value (expected '0' or '1')"), deadline)
		zmq.Destroy()
	}()

	for {
		err := ProxyReqRep1(ws, zmq)
		if err != nil {
			log.Printf("Error in forwarding request from websocket->zeromq: %s", err)
			break
		}
		err = ProxyReqRep2(ws, zmq)
		if err != nil {
			log.Printf("Error in forwarding request from zeromq->websocket: %s", err)
			break
		}
	}
}

func ProxyReqRep1(ws *websocket.Conn, zmq *czmq.Sock) error {
	done := false

	for !done {
		messageType, message, err := ws.ReadMessage()
		log.Printf("ReqRep1: received %d: %s\n", messageType, string(message))
		if err != nil {
			log.Printf("Error conn.ReadMessage: %s", err)
			return err
		}
		if messageType != websocket.TextMessage {
			log.Printf("Got unexpected message type: %s", messageType)
			return &InvalidMessageTypeError{messageType}
		}

		flags := czmq.FlagNone
		switch more := message[0]; more {
		case FINAL_FRAME:
			// Final frame, let's break out of the loop afterthis iteration
			done = true
		case MORE_FRAME: // There are more frames after this one.
			flags |= czmq.FlagMore
		default:
			log.Printf("Got an invalid MORE frame value: %c", more)
			return &InvalidMoreFlagError{more}
		}

		log.Printf("Forwarding 1 frame: final(%s) - %s", message[0] == FINAL_FRAME, string(message[1:]))
		err = zmq.SendFrame(message[1:], flags)
		if err != nil {
			log.Printf("zmq.SendFrame failed: %s\n", err)
			return err
		}
	}
	return nil
}

func ProxyReqRep2(ws *websocket.Conn, zmq *czmq.Sock) error {
	done := false
	for !done {
		frame, more, err := zmq.RecvFrame()
		if err != nil {
			log.Printf("zmq.RecvFrame failed: %s", err)
			return err
		}

		log.Printf("ReqRep2: received %s: %s\n", more, string(frame))

		var payload [1]byte
		if more&czmq.FlagMore == czmq.FlagMore {
			payload[0] = MORE_FRAME
		} else {
			payload[0] = FINAL_FRAME
			done = true
		}
		err = ws.WriteMessage(websocket.TextMessage, append(payload[:], frame...))
		if err != nil {
			log.Printf("ws.WriteMessage failed: %s", err)
			return err
		}
	}
	return nil
}

func ProxySub(ws *websocket.Conn, zmq *czmq.Sock) {
	log.Printf("Proxy SUB")
	defer func() {
		deadline := time.Now().Add(1 * time.Second)
		ws.WriteControl(websocket.CloseMessage, []byte("Invalid MORE frame value (expected '0' or '1')"), deadline)
		zmq.Destroy()
	}()

	for {
    frame, more, err := zmq.RecvFrame()
		var payload [1]byte
		if more&czmq.FlagMore == czmq.FlagMore {
			payload[0] = MORE_FRAME
		} else {
			payload[0] = FINAL_FRAME
		}
    err = ws.WriteMessage(websocket.TextMessage, append(payload[:], frame...))
    if err != nil {
      log.Printf("ws.WriteMessage failed: %s", err)
      break
    }
	}
}

