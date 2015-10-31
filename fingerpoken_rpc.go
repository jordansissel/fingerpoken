package main

import (
	"golang.org/x/net/websocket"
	"net/rpc/jsonrpc"
	"os/exec"
)

func ws2jsonrpc(ws *websocket.Conn) {
	jsonrpc.ServeConn(ws)
}

//func (t *T) MethodName(argType T1, replyType *T2) error

type LIRC struct{}

type LIRCResponse struct {
	Success bool
	Message string
}

type RemoteCode struct {
	Remote string
	Codes  []string
}

func (l *LIRC) send_once(rc *RemoteCode, reply *LIRCResponse) {
	var args []string = []string{"irsend"}
	args = append(args, rc.Remote)
	args = append(args, rc.Codes...)
	cmd := exec.Command("irsend", args...)
	output, err := cmd.Output()

	reply.Success = err != nil
	reply.Message = string(output)
}
