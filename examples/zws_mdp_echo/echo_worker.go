package main

import (
  "log"
)

type echoWorker struct{}

func (e *echoWorker) Request(request [][]byte) (response [][]byte, err error) {
	for i, x := range request { log.Printf("Echo worker: frame %d: %v (%s)\n", i, x, string(x)) }
	response = append(response, request...)
	return
}

func (e *echoWorker) Heartbeat()  {}
func (e *echoWorker) Disconnect() {}
