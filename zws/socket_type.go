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
	"fmt"
	czmq "github.com/zeromq/goczmq"
	"log"
)

type SocketType int8

const (
	INVALID = iota
	REQ     = iota
	DEALER  = iota
	// TODO(sissel): Fill in more as we support them.
)

func ParseSocketType(name string) (SocketType, error) {
	switch name {
	case "req":
		return REQ, nil
	case "dealer":
		return DEALER, nil
	}
	return INVALID, fmt.Errorf("Invalid socket type: %s", name)
}

func (s *SocketType) String() string {
	switch *s {
	case REQ:
		return "<REQ>"
	case DEALER:
		return "<DEALER>"
	}
	return "<INVALID_SOCKET 0x%x>"
}

func (s *SocketType) EndpointSuffix() (string, error) {
	switch *s {
	case REQ:
		return "req", nil
	case DEALER:
		return "dealer", nil
	}
	return "", &InvalidSocketTypeError{*s}
}

func (s *SocketType) isValid() bool {
	switch *s {
	case REQ, DEALER:
		return true
	}

	return false
}

func (s *SocketType) Create(endpoint string) (sock *czmq.Sock, err error) {
	log.Printf("New `%s` on `%s`", s, endpoint)
	switch *s {
	case REQ:
		return czmq.NewReq(endpoint)
	case DEALER:
		return czmq.NewDealer(endpoint)
	}

	return nil, &InvalidSocketTypeError{*s}
}
