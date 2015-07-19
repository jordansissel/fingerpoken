package mdp

// Frame 1: "MDPW01" (six bytes, representing MDP/Worker v0.1)
var MDP_CLIENT = []byte("MDPC01")
var MDP_WORKER = []byte("MDPW01")

// Commands
const C_READY = byte(0x01)
const C_REQUEST = byte(0x02)
const C_REPLY = byte(0x03)
const C_HEARTBEAT = byte(0x04)
const C_DISCONNECT = byte(0x05)

