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
	// TODO(sissel): Fill in more as we support them.
)

func ParseSocketType(name string) (SocketType, error) {
	switch name {
	case "req":
		return REQ, nil
	}
	return INVALID, fmt.Errorf("Invalid socket type: %s", name)
}

func (s *SocketType) String() string {
	switch *s {
	case REQ:
		return "<REQ>"
	}
	return "<INVALID_SOCKET 0x%x>"
}

func (s *SocketType) EndpointSuffix() string {
	switch *s {
	case REQ:
		return "req"
	}
	panic("Invalid socket type")
}

func (s *SocketType) isValid() bool {
	switch *s {
	case REQ:
		return true
	}

	return false
}

func (s *SocketType) Create(endpoint string) (sock *czmq.Sock, err error) {
	log.Printf("New `%s` on `%s`", s, endpoint)
	switch *s {
	case REQ:
		return czmq.NewReq(endpoint)
	}

	log.Fatalf("Invalid socket type given: %s. You should check isValid() before calling this.", s)
	return
}
