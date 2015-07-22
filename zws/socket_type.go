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
