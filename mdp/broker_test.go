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
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/jordansissel/fingerpoken/zap"
	czmq "github.com/jordansissel/goczmq"
	"testing"
	"time"
)

func TestBrokerWorkerIntegration(t *testing.T) {
	broker_address := fmt.Sprintf("inproc://%s", randomHex())
	service := randomHex()
	w := NewWorker(broker_address, service)
	b, _ := NewBroker(broker_address)
	go b.Run()
	err := w.ensure_connected()
	if err != nil {
		t.Errorf("Worker#ensure_connected failed: %s\n", err)
		return
	}
}

func TestEndToEnd(t *testing.T) {
	broker_address := fmt.Sprintf("inproc://%s", randomHex())
	service := randomHex()
	w := NewWorker(broker_address, service)
	c := NewClient(broker_address)
	b, err := NewBroker(broker_address)
	if err != nil {
		t.Errorf("NewBroker(%v) failed: %s", broker_address, err)
		return
	}

	go b.Run()
	go w.Run(&helloGreeter{})

	text := randomHex()
	body := [1][]byte{[]byte(text)}
	// TODO(sissel): Send a client request, verify response
	response, err := c.SendRecv(service, body[:])
	//for i, x := range response { log.Printf("Client sending): frame %d: %v (%s)\n", i, x, string(x)) }
	if err != nil {
		t.Errorf("Failure in request to service `%s`: %s\n", service, err)
		return
	}

	if expected, actual := 1, len(response); expected != actual {
		t.Errorf("Expected %d frames for reply. Got %d frames\n", expected, actual)
		return
	}
	if !bytes.Equal(response[0], helloGreeting) {
		t.Errorf("Response did not match `%s`", string(helloGreeting))
		return
	}
}

type HeartbeatChecker struct {
	callback func(time.Time)
}

func (h *HeartbeatChecker) Request(request [][]byte) (response [][]byte, err error) {
	return
}
func (h *HeartbeatChecker) Heartbeat() {
	h.callback(time.Now())
}
func (h *HeartbeatChecker) Disconnect() {}

func TestBrokerToWorkerHeartbeat(t *testing.T) {
	broker_address := fmt.Sprintf("inproc://%s", randomHex())
	service := randomHex()
	w := NewWorker(broker_address, service)
	b, err := NewBroker(broker_address)
	if err != nil {
		t.Errorf("NewBroker(%v) failed: %s", broker_address, err)
		return
	}

	b.HeartbeatInterval = 10 * time.Millisecond

	beats := make(chan time.Time)
	hc := &HeartbeatChecker{
		callback: func(now time.Time) {
			log.Printf("Heartbeat!")
			beats <- now
		},
	}
	start := time.Now()
	go w.Run(hc)
	go b.Run()

	heartbeatTime := <-beats
	if !heartbeatTime.After(start) {
		t.Errorf("Heartbeat time occurred in the past?!")
		return
	}
}

func TestWorkerToBrokerHeartbeat(t *testing.T) {
	broker_address := fmt.Sprintf("inproc://%s", randomHex())
	service := randomHex()
	w := NewWorker(broker_address, service)
	b, err := NewBroker(broker_address)
	if err != nil {
		t.Errorf("NewBroker(%v) failed: %s", broker_address, err)
		return
	}

	b.HeartbeatInterval = 10 * time.Millisecond
	w.HeartbeatInterval = b.HeartbeatInterval

	beats := make(chan time.Time)
	b.HeartbeatCallback = func(entry *workerEntry) {
		beats <- time.Now()
	}
	start := time.Now()
	go w.Run(&helloGreeter{})
	go b.Run()

	heartbeatTime := <-beats
	if !heartbeatTime.After(start) {
		t.Errorf("Heartbeat time occurred in the past?!")
		return
	}
}

func TestEndToEndWithCurve(t *testing.T) {
	agent, _ := zap.NewZapAgent()
	go agent.Run(&zap.OpenAccess{})
	defer agent.Destroy()

	broker_address := fmt.Sprintf("inproc://%s", randomHex())
	service := randomHex()
	b, err := NewBroker(broker_address)
	if err != nil {
		t.Errorf("NewBroker(%v) failed: %s", broker_address, err)
		return
	}

	server_cert := czmq.NewCert()
	b.CurveCertificate = server_cert

	port, err := b.Bind("tcp://127.0.0.1:*")
	if err != nil {
		t.Errorf("Broker.Bind failed: %s", err)
		return
	}

	broker_local_address := fmt.Sprintf("tcp://127.0.0.1:%d", port)

	w := NewWorker(broker_local_address, service)
	w.CurveServerPublicKey = server_cert.PublicText()

	c := NewClient(broker_local_address)
	defer c.Destroy()
	c.CurveServerPublicKey = server_cert.PublicText()

	go b.Run()
	go w.Run(&helloGreeter{})

	text := randomHex()
	body := [1][]byte{[]byte(text)}
	// TODO(sissel): Send a client request, verify response
	response, err := c.SendRecv(service, body[:])
	//for i, x := range response { log.Printf("Client sending): frame %d: %v (%s)\n", i, x, string(x)) }
	if err != nil {
		t.Errorf("Failure in request to service `%s`: %s\n", service, err)
		return
	}

	if expected, actual := 1, len(response); expected != actual {
		t.Errorf("Expected %d frames for reply. Got %d frames\n", expected, actual)
		return
	}
	if !bytes.Equal(response[0], helloGreeting) {
		t.Errorf("Response did not match `%s`", string(helloGreeting))
		return
	}
}
