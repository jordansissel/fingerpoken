package mdp

import (
	"fmt"
)

// Frame 1: "MDPW01" (six bytes, representing MDP/Worker v0.1)
var mdp_CLIENT = []byte("MDPC01")
var mdp_WORKER = []byte("MDPW01")

// commands
type command byte

const (
	_               = iota
	c_READY command = iota
	c_REQUEST
	c_REPLY
	c_HEARTBEAT
	c_DISCONNECT
)

func (c command) String() string {
	switch c {
	case c_READY:
		return "<READY>"
	case c_REQUEST:
		return "<REQUEST>"
	case c_REPLY:
		return "<REPLY>"
	case c_HEARTBEAT:
		return "<HEARTBEAT>"
	case c_DISCONNECT:
		return "<DISCONNECT>"
	default:
		return fmt.Sprintf("<INVALID COMMAND 0x%0x>", byte(c))
	}
}

// commands that never change
var m_HEARTBEAT = [3][]byte{
	[]byte{},                  // SPEC: Frame 0: Empty frame
	mdp_WORKER,                // SPEC: Frame 1: "MDPW01" (six bytes, representing MDP/Worker v0.1)
	[]byte{byte(c_HEARTBEAT)}, // SPEC: Frame 2: 0x04 (one byte, representing HEARTBEAT)
}

var m_DISCONNECT = [3][]byte{
	[]byte{},                   // SPEC: Frame 0: Empty frame
	mdp_WORKER,                 // SPEC: Frame 1: "MDPW01" (six bytes, representing MDP/Worker v0.1)
	[]byte{byte(c_DISCONNECT)}, // SPEC: Frame 2: 0x04 (one byte, representing HEARTBEAT)
}
