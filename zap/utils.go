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
	czmq "github.com/jordansissel/goczmq"
)

func czmqSockSafeRecv(sock *czmq.Sock) (frame []byte, more int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Panic while calling Recv(): %s", r)
			return
		}
	}()
	frame, more, err = sock.RecvFrame()
	return
}

func czmqSockSafeRecvMessage(sock *czmq.Sock) (frames [][]byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Panic while calling RecvFrames(): %s", r)
			return
		}
	}()
	frames, err = sock.RecvMessage()
	return
}

type OpenAccess struct{}

func logAuthRequest(authRequest ZapRequest) {
	log.Printf("Auth[%s] from %s", authRequest.Mechanism, authRequest.Address)
	switch authRequest.Mechanism {
	case Curve:
		log.Printf("Client public key: %v", authRequest.Credentials[0])
	case Plain:
		log.Printf("Client user: %s", authRequest.Credentials[0])
	}
}

func (o *OpenAccess) Authorize(authRequest ZapRequest) (status Status, err error) {
	logAuthRequest(authRequest)
	return Success, nil
}

var NULL_PRIVATE_KEY [32]byte
type DenyAccess struct{}

func (d *DenyAccess) Authorize(authRequest ZapRequest) (status Status, err error) {
	logAuthRequest(authRequest)
	return AuthenticationFailure, nil
}

type RestrictedAccess struct {
	//Store *czmq.CertStore
	Trusted map[string]bool
}

func NewRestrictedAccess() (r *RestrictedAccess) {
  r = &RestrictedAccess{
    Trusted: make(map[string]bool),
  }
  return r
}

func (r *RestrictedAccess) Allow(public_key string) {
  //r.Store.Insert(cert)
  r.Trusted[public_key] = true
}

func (r *RestrictedAccess) Authorize(authRequest ZapRequest) (status Status, err error) {
	if authRequest.Mechanism != "CURVE" {
		return AuthenticationFailure, fmt.Errorf("Auth mechanism must be CURVE")
	}
	if len(authRequest.Credentials) != 1 {
		return AuthenticationFailure, fmt.Errorf("Invalid credentials format")
	}

	client_key := string(authRequest.Credentials[0])
	if !r.Trusted[client_key] {
    log.Info("Trusting this client for next time")
    r.Allow(client_key)
		return AuthenticationFailure, nil
	}
	return Success, nil
}
