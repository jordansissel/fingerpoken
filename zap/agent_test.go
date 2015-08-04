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
package zap

import (
	"bytes"
	"fmt"
	czmq "github.com/jordansissel/goczmq"
	"math/rand"
	"testing"
)

func randomHex() (value string) {
	var length = rand.Int31n(10) + 5
	for x := int32(0); x < length; x += 1 {
		value = fmt.Sprintf("%s%02x", value, rand.Int31n(256))
	}
	return
}

func TestZAPAllow(t *testing.T) {
	agent, _ := NewZapAgent()
	go agent.Run(&OpenAccess{})
	defer agent.Destroy()

	domain := randomHex()
	requestId := randomHex()
	ip := "1.2.3.4"
	id := randomHex()
	mechanism := randomHex()
	credentials := randomHex()

	status, err := sendAuthRequest(requestId, domain, ip, id, mechanism, credentials)

	if status != Success {
		t.Errorf("Expected Success, got %s. %s", status, err)
		return
	}
}

func TestZAPDeny(t *testing.T) {
	agent, _ := NewZapAgent()
	go agent.Run(&DenyAccess{})
	defer agent.Destroy()

	domain := randomHex()
	requestId := randomHex()
	ip := "1.2.3.4"
	id := randomHex()
	mechanism := randomHex()
	credentials := randomHex()

	status, err := sendAuthRequest(requestId, domain, ip, id, mechanism, credentials)

	if status != AuthenticationFailure {
		t.Errorf("Expected Authentication Failure, got %s. %s", status, err)
		return
	}
}

func sendAuthRequest(requestId, domain, ip, id, mechanism, credentials string) (Status, error) {
	sock, err := czmq.NewReq(ZAP_ENDPOINT)
	if err != nil {
		return InternalError, fmt.Errorf("Error creating new REQ for %s: %s", ZAP_ENDPOINT, err)
	}

	err = sock.SendMessage([][]byte{
		[]byte(ZAP_VERSION), //The version frame, which SHALL contain the three octets "1.0".
		[]byte(requestId),   //The request id, which MAY contain an opaque binary blob.
		[]byte(domain),      //The domain, which SHALL contain a string.
		[]byte(ip),          //The address, the origin network IP address.
		[]byte(id),          //The identity, the connection Identity, if any.
		[]byte(mechanism),   //The mechanism, which SHALL contain a string.
		[]byte(credentials), //The credentials, which SHALL be zero or more opaque frames.
	})
	if err != nil {
		return InternalError, fmt.Errorf("Error in SendMessage: %s", err)
	}

	reply, err := sock.RecvMessage()
	//for i, x := range reply { log.Printf("reply %d: %s", i, string(x)) }

	if err != nil {
		return InternalError, fmt.Errorf("Error in SendMessage: %s", err)
	}

	// SPEC: The version frame, which SHALL contain the three octets "1.0".
	if !bytes.Equal(reply[0], []byte(ZAP_VERSION)) {
		return InternalError, fmt.Errorf("Expected first frame to be `%s`.", ZAP_VERSION)
	}

	// SPEC: The request id, which MAY contain an opaque binary blob.
	if !bytes.Equal(reply[1], []byte(requestId)) {
		return InternalError, fmt.Errorf("Request ID did not match")
	}

	// SPEC: The status code, which SHALL contain a string.
	status := Status(reply[2])
	reason := string(reply[3])
	switch status {
	case Success:
		return status, nil
	case TemporaryError:
		return status, fmt.Errorf("Temporary Error: %s", reason)
	case AuthenticationFailure:
		return status, fmt.Errorf("Authentication Failure: %s", reason)
	case InternalError:
		return status, fmt.Errorf("Internal Server Error: %s", reason)
	default:
		return status, fmt.Errorf("Invalid status code: %s", status)
	}

	// SPEC: The status text, which MAY contain a string.
	// SPEC: The user id, which SHALL contain a string.
	// SPEC: The metadata, which MAY contain a blob.
}
