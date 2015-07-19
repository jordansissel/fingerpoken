package mdp

import (
	"fmt"
)

// Frame 1: "MDPW01" (six bytes, representing MDP/Worker v0.1)
var MDP_CLIENT = []byte("MDPC01")
var MDP_WORKER = []byte("MDPW01")

// Commands
type Command byte

const (
	_               = iota
	C_READY Command = iota
	C_REQUEST
	C_REPLY
	C_HEARTBEAT
	C_DISCONNECT
)

func (c Command) String() string {
	switch c {
	case C_READY:
		return "<READY>"
	case C_REQUEST:
		return "<REQUEST>"
	case C_REPLY:
		return "<REPLY>"
	case C_HEARTBEAT:
		return "<HEARTBEAT>"
	case C_DISCONNECT:
		return "<DISCONNECT>"
	default:
		return fmt.Sprintf("<INVALID COMMAND 0x%0x>", byte(c))
	}
}

// Commands that never change
var M_HEARTBEAT = [3][]byte{
	[]byte{},                  // SPEC: Frame 0: Empty frame
	MDP_WORKER,                // SPEC: Frame 1: "MDPW01" (six bytes, representing MDP/Worker v0.1)
	[]byte{byte(C_HEARTBEAT)}, // SPEC: Frame 2: 0x04 (one byte, representing HEARTBEAT)
}

var M_DISCONNECT = [3][]byte{
	[]byte{},                   // SPEC: Frame 0: Empty frame
	MDP_WORKER,                 // SPEC: Frame 1: "MDPW01" (six bytes, representing MDP/Worker v0.1)
	[]byte{byte(C_DISCONNECT)}, // SPEC: Frame 2: 0x04 (one byte, representing HEARTBEAT)
}
