package mdp

import (
	"bytes"
	"fmt"
	"github.com/jordansissel/fingerpoken/util"
	"io/ioutil"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type JSONRPCWorker struct {
	rpc    *rpc.Server
	worker *Worker
}

func NewJSONRPCWorker(broker_endpoint, service string) (j *JSONRPCWorker) {
	j = &JSONRPCWorker{}
	j.rpc = rpc.NewServer()
	j.worker = NewWorker(broker_endpoint, service)
	return j
}

func (j *JSONRPCWorker) Register(handler interface{}) error {
	return j.rpc.Register(handler)
}

func (j *JSONRPCWorker) Run() {
	j.worker.Run(j)
}

// the mdp.Worker interface
func (j *JSONRPCWorker) Request(request [][]byte) (response [][]byte, err error) {
	if len(request) != 1 {
		// TODO(sissel): invalid request. Must be one frame only.
		log.Printf("Received JSONRPC Request w/ %d frames, requires exactly 1 frame.", len(request))
		return nil, fmt.Errorf("Received JSONRPC Request w/ %d frames, requires exactly 1 frame.", len(request))
	}
	//log.Printf("JSONRPC: %#v\n", string(request[0]))
	input := bytes.NewBuffer(request[0])
	output := &bytes.Buffer{}
	codec := jsonrpc.NewServerCodec(util.Buffer{Reader: input, Writer: output, Closer: ioutil.NopCloser(nil)})
	err = j.rpc.ServeRequest(codec)
	if err != nil {
		// Something wrong with this rpc call.
		log.Printf("rpc.ServeRequest failed: %s", err)
		return nil, fmt.Errorf("rpc.ServeRequest failed: %s", err)
	}
	response = append(response, output.Bytes())
	return
}

func (j *JSONRPCWorker) Heartbeat()  {}
func (j *JSONRPCWorker) Disconnect() {}
