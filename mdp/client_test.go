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
	payload := [1][]byte{[]byte(randomHex())}

	router, err := czmq.NewRouter(endpoint)
	defer router.Destroy()
	if err != nil {
		t.Errorf("Creating new router failed, %s: %s", endpoint, err)
		return
	}

	client := NewClient(endpoint)
	defer client.Destroy()
	go client.SendRecv(service, payload[:])

	frames, err := router.RecvMessage()
	if err != nil {
		t.Errorf("Error while reading from router %s: %s", endpoint, err)
		return
	}

	if count := len(frames); count < 5 {
		t.Errorf("Majordomo requests must have at least 5 frames, got %d frames\n", count)
		return
	}

	//for i, x := range frames { fmt.Printf("%d: %v (%s) \n", i, x, string(x)) }

	// frames[0] is the client/session id for the router socket, ignore it.
	// frames[1 ... ] are the actual request

	// Frame 0: Empty (zero bytes, invisible to REQ application)
	if len(frames[1]) != 0 {
		t.Errorf("Majordomo request frame #1 must be empty.\n")
		return
	}

	if !bytes.Equal(frames[2], mdp_CLIENT) {
		t.Errorf("Majordomo request frame #2 must be `%s`, got `%s`\n", string(mdp_CLIENT), string(frames[1]))
		return
	}

	//if len(frames[3]) != 1 || command(frames[3][0]) != c_REQUEST {
	//t.Errorf("Majordomo request frame #3 must be REQUEST")
	//return
	//}

	if !bytes.Equal(frames[3], []byte(service)) {
		t.Errorf("Majordomo request frame #4 must be a service name. Expected `%s`, got `%s`", service, string(frames[4]))
		return
	}

	if expected, actual := len(payload), len(frames[4:]); expected != actual {
		t.Errorf("Expected body with %d frames, got %d frames\n", expected, actual)
		return
	}

	if !bytes.Equal(frames[4], payload[0]) {
		t.Errorf("Majordomo request body did not match.")
		return
	}
}
