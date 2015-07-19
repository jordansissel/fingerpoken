package mdp

import (
	"bytes"
	"fmt"
	czmq "github.com/zeromq/goczmq"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient("inproc://foo")
	err := c.ensure_connected()
	if err != nil {
		t.Errorf("Client#ensure_connected failed: %s\n", err)
	}
}

func TestNewClientWithInvalidEndpoint(t *testing.T) {
	c := NewClient("nonsense")
	err := c.ensure_connected()
	if err == nil {
		t.Errorf("Client#ensure_connected shoudl have failed.")
	}
}

func TestClientSendFraming(t *testing.T) {
	// Randomize things for better testing confidence
	endpoint := fmt.Sprintf("inproc://%s", randomHex())
	service := randomHex()
	payload := randomHex()

	router, err := czmq.NewRouter(endpoint)
	if err != nil {
		t.Errorf("Creating new router failed, %s: %s", endpoint, err)
		return
	}

	client := NewClient(endpoint)
	go func() {
		//response, err := client.Send("fancy", []byte("hello world"))
		client.Send(service, []byte(payload))
	}()

	frames, err := router.RecvMessage()
	if err != nil {
		t.Errorf("Error while reading from router %s: %s", endpoint, err)
		return
	}

	if count := len(frames); count < 5 {
		t.Errorf("Majordomo requests must have at least 5 frames, got %d frames\n", count)
		return
	}

	//for i, x := range frames {
	//fmt.Printf("%d: %v (%s) \n", i, x, string(x))
	//}

	// frames[0] is the client/session id for the router socket, ignore it.
	// frames[1 ... ] are the actual request

	// Frame 0: Empty (zero bytes, invisible to REQ application)
	if len(frames[1]) != 0 {
		t.Errorf("Majordomo request frame #1 must be empty.\n")
		return
	}

	if !bytes.Equal(frames[2], MDP_CLIENT) {
		t.Errorf("Majordomo request frame #2 must be `%s`, got `%s`\n", string(MDP_CLIENT), string(frames[1]))
		return
	}

	if !bytes.Equal(frames[3], []byte(service)) {
		t.Errorf("Majordomo request frame #3 must be a service name. Expected `%s`, got `%s`", service, string(frames[2]))
		return
	}

	if !bytes.Equal(frames[4], []byte(payload)) {
		t.Errorf("Majordomo request frame #4 must be a service name. Expected `%s`, got `%s`", payload, string(frames[3]))
		return
	}
}
