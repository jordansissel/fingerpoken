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
package mdp

import (
	"fmt"
	czmq "github.com/zeromq/goczmq"
	"log"
	"time"
)

// TODO(sissel): turn this into an interface?
type Worker struct {
	sock                *czmq.Sock
	broker              string
	service             string
	HeartbeatInterval   time.Duration
	MaxMissedHeartbeats int64

	brokerExpiration time.Time
	poller           *czmq.Poller
}

func NewWorker(broker string, service string) (w *Worker) {
	w = &Worker{
		broker:              broker,
		service:             service,
		HeartbeatInterval:   5 * time.Second,
		MaxMissedHeartbeats: 3,
	}
	return
}

type RequestHandler interface {
	Request([][]byte) ([][]byte, error)
	Heartbeat()
	Disconnect()
}

func (w *Worker) Run(requestHandler RequestHandler) error {
	w.ensure_connected()
	nextHeartbeat := time.Now().Add(w.HeartbeatInterval)
	w.brokerExpiration = w.nextExpiration()
	for {
		s := czmqPollerSafeWait(w.poller, durationInMilliseconds(w.HeartbeatInterval))
		if s != nil {
			// Data is ready, let's process something.
			client, cmd, body, err := w.readRequest()
			if err != nil {
				log.Printf("Worker: Error reading request: %s\n", err)
				w.Reset()
				continue
			}
			err = w.handleCommand(requestHandler, client, cmd, body)
			if err != nil {
				log.Printf("Error processing cmd %s", cmd)
				w.Reset()
			}
		}

		now := time.Now()
		if now.After(w.brokerExpiration) {
			// Broker hasn't been heard from in a while. Let's reset.
			log.Printf("Broker hasn't been heard from in %s. Resetting connection.", now.Sub(w.brokerExpiration))
			w.Reset()
		}
		if now.After(nextHeartbeat) {
			// It's time to send a heartbeat to the broker.
			w.sendHeartbeat()
			nextHeartbeat = time.Now().Add(w.HeartbeatInterval)
		}
	}
}

func (w *Worker) handleCommand(requestHandler RequestHandler, client []byte, cmd command, body [][]byte) error {
	// Got a cmd from the broker, let's update the expiration.
	w.brokerExpiration = w.nextExpiration()

	switch cmd {
	default:
		log.Printf("Worker: Got an invalid cmd from broker. Will reset connection. (cmd: %v)", cmd)
		w.Reset()
	case c_HEARTBEAT:
		requestHandler.Heartbeat()
	case c_DISCONNECT:
    log.Printf("Worker: Received DISCONNECT from broker. Will reset connection")
		requestHandler.Disconnect()
		w.Reset()
	case c_REQUEST:
		// The spec supports multiple frames for a request message. Let's support that.
		reply_body, err := requestHandler.Request(body)
		if err != nil {
			log.Printf("Worker: Error handling request: %s\n", err)
			return err
		}

		frames := append([][]byte{
			[]byte{},
			mdp_WORKER,
			[]byte{byte(c_REPLY)},
			client,
			[]byte{}, // SPEC Frame 4: Empty (zero bytes, envelope delimiter)
		}, reply_body...)

		//for i, x := range frames { log.Printf("Worker(to Broker): frame %d: %v (%s)\n", i, x, string(x)) }
		err = w.sock.SendMessage(frames)
		if err != nil {
			log.Printf("Worker: Error sending reply: %s\n", err)
			return err
		}
	} // switch cmd

	return nil
}

func (w *Worker) readRequest() (client []byte, cmd command, body [][]byte, err error) {
	frames, err := w.sock.RecvMessage()
	if err != nil {
		return
	}
	//for i, x := range frames { log.Printf("Worker(via Broker): frame %d: %v (%s)\n", i, x, string(x)) }

	err = validateWorkerHeader(frames[:])
	if err != nil {
		err = fmt.Errorf("Got an invalid worker header in request. Will reset connection. %s", err)
		return
	}

	cmd = command(frames[2][0])
	if cmd == c_REQUEST {
		if len(frames) < 5 {
			err = fmt.Errorf("Got a request with too-few frames. Expected at least 5, got %d", len(frames))
			return
		}
		client = frames[3]
		body = frames[5:]
		//for i, x := range body { log.Printf("Worker(request body): frame %d: %v (%s)\n", i, x, string(x)) }
	}
	return
}

func (w *Worker) Reset() {
	w.sock.Destroy()
	w.sock = nil
	w.poller.Destroy()
	w.poller = nil
	w.ensure_connected()
}

func (w *Worker) ensure_connected() error {
	if w.sock != nil {
		return nil
	}

	var err error
	w.sock, err = czmq.NewDealer(w.broker)
	if err != nil {
		return err
	}
	w.poller, err = czmq.NewPoller(w.sock)
	if err != nil {
		w.sock.Destroy()
		return err
	}

	err = w.sendReady()
	if err != nil {
		return err
	}
	return nil
}

func (w *Worker) sendReady() (err error) {
	message := [4][]byte{
		[]byte{},              // SPEC: Frame 0: Empty frame
		mdp_WORKER,            // SPEC: Frame 1: "MDPW01" (six bytes, representing MDP/Worker v0.1)
		[]byte{byte(c_READY)}, // SPEC: Frame 2: 0x01 (one byte, representing READY)
		[]byte(w.service),     // SPEC: Frame 3: Service name (printable string)
	}

	err = w.sock.SendMessage(message[:])
	// SPEC: There is no response to a READY. The worker SHOULD assume the registration succeeded.
	if err != nil {
		return err
	}
	return nil
}

func (w *Worker) sendHeartbeat() (err error) {
	//log.Printf("Worker: Sending heartbeat")
	err = w.sock.SendMessage(m_HEARTBEAT[:])
	if err != nil {
		return err
	}
	return nil
}

func (w *Worker) sendDisconnect() (err error) {
	err = w.sock.SendMessage(m_DISCONNECT[:])
	if err != nil {
		return err
	}
	return nil
}

func (w *Worker) nextExpiration() time.Time {
	// SPEC: A peer MUST consider the other peer "disconnected" if no heartbeat arrives within some multiple of that interval (usually 3-5).
	return time.Now().Add(time.Duration(int64(w.HeartbeatInterval) * w.MaxMissedHeartbeats))
}
