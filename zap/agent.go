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
	"fmt"
	log "github.com/Sirupsen/logrus"
	czmq "github.com/zeromq/goczmq"
)

type InvalidZAPRequest struct {
	detail string
}

func (e *InvalidZAPRequest) Error() string {
	return fmt.Sprintf("Invalid ZAP request: %s", e.detail)
}

type ZapRequest struct {
	Domain      string
	Address     string
	Identity    []byte
	Mechanism   Mechanism
	Credentials [][]byte
}

type ZapAgent struct {
	zmq     *czmq.Sock
	handler ZapHandler
}

type ZapHandler interface {
	Authorize(authRequest ZapRequest) (Status, error)
}

const ZAP_ENDPOINT = "inproc://zeromq.zap.01"
const ZAP_VERSION = "1.0"

type Status string
type Mechanism string

const (
	Success               Status = "200"
	TemporaryError        Status = "300"
	AuthenticationFailure Status = "400"
	InternalError         Status = "500"

	Curve Mechanism = "CURVE"
	Null  Mechanism = "NULL"
	Plain Mechanism = "PLAIN"
)

func NewZapAgent() (zap *ZapAgent, err error) {
	zap = &ZapAgent{}
	zap.zmq, err = czmq.NewRouter(ZAP_ENDPOINT)
	return
}

func (zap *ZapAgent) SetHandler(handler ZapHandler) {
	zap.handler = handler
}

func (zap *ZapAgent) Destroy() {
	zap.zmq.Destroy()
	zap.zmq = nil
}

func (zap *ZapAgent) Run(handler ZapHandler) error {
	// SPEC: An address delimiter frame, which SHALL have a length of zero.
	// SPEC: The version frame, which SHALL contain the three octets "1.0".
	// SPEC: The request id, which MAY contain an opaque binary blob.
	// SPEC: The domain, which SHALL contain a string.
	// SPEC: The address, the origin network IP address.
	// SPEC: The identity, the connection Identity, if any.
	// SPEC: The mechanism, which SHALL contain a string.
	// SPEC: The credentials, which SHALL be zero or more opaque frames.
	zap.SetHandler(handler)
	for {
		err := zap.Once(zap.handler)
		if err != nil {
			return err
		}
	}

	return nil
}

func (zap *ZapAgent) Once(handler ZapHandler) error {
	frames, err := czmqSockSafeRecvMessage(zap.zmq)
	if err != nil {
		return err
	}

	replyTo := frames[0]
	message := frames[1:]

	if err := validateRequest(message); err != nil {
		return err
	}

	requestId := message[2]
	authRequest := ZapRequest{
		Domain:    string(message[3]),
		Address:   string(message[4]),
		Identity:  message[5],
		Mechanism: Mechanism(message[6]),
	}

	if len(message) > 7 {
		authRequest.Credentials = message[7:]
	}

	var status Status
	//log.Printf("Auth[%s] from %s", authRequest.Mechanism, authRequest.Address)
	if handler == nil {
		// Default deny.
		log.Printf("No ZAP handler was given, so I cannot authorize anything. Denying %S from %s.", authRequest.Mechanism, authRequest.Address)
		status = InternalError
		err = fmt.Errorf("No ZAP handler available")
	} else {
		status, err = handler.Authorize(authRequest)
	}

	statusText := ""
	if status != Success && err != nil {
		// TODO(sissel): differentiate between temporary and internal errors.
		statusText = err.Error()
	}

	response := [][]byte{
		replyTo,
		[]byte{},      // address delimiter
		[]byte("1.0"), // version
		requestId,
		[]byte(string(status)), // success? failure?
		[]byte(statusText),     // details
		[]byte{},               // user id. Optional?
		[]byte{},               // metadata, optional?
	}
	log.Printf("%s auth from %s", status, authRequest.Address)

	err = zap.zmq.SendMessage(response)
	return err
}

func validateRequest(message [][]byte) error {
	if len(message) < 7 {
		return &InvalidZAPRequest{detail: fmt.Sprintf("ZAP requests have at least 7 frames, got %d frames", len(message))}
	}

	if len(message[0]) != 0 {
		return &InvalidZAPRequest{detail: "First frame must be length 0"}
	}

	if string(message[1]) != "1.0" {
		return &InvalidZAPRequest{detail: "Second frame (version) must '1.0'"}
	}

	// message[2] is request id, can be anything
	// message[3] is domain, must be a string, but I think can be anything
	// message[4] must be an address, probably nonzero length.

	if len(message[4]) == 0 {
		return &InvalidZAPRequest{detail: "Fifth frame (address) must not be empty"}
	}

	// message[5] is identity, optional
	// message[6] is the mechanism

	if len(message[6]) == 0 {
		return &InvalidZAPRequest{detail: "Seventh frame (mechanism) must not be empty"}
	}

	return nil
}
