package mdp

import (
	"bytes"
	"fmt"
	"testing"
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
	go w.Run(&HelloGreeter{})

	text := randomHex()
	body := [1][]byte{[]byte(text)}
	// TODO(sissel): Send a client request, verify response
	response, err := c.Send(service, body[:])
	if err != nil {
		t.Errorf("Failure in request to service `%s`: %s\n", service, err)
		return
	}

	if expected, actual := 1, len(response); expected != actual {
		t.Errorf("Expected %d frames for reply. Got %d frames\n", expected, actual)
		return
	}
	if !bytes.Equal(response[0], HELLO_GREETING) {
		t.Errorf("Response did not match `%s`", string(HELLO_GREETING))
		return
	}
}
