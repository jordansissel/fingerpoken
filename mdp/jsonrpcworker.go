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
package mdp

import (
	"bytes"
	"fmt"
	"github.com/jordansissel/fingerpoken/util"
	czmq "github.com/zeromq/goczmq"
	"io/ioutil"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
	"reflect"
)

type JSONRPCWorker struct {
	CurveServerPublicKey string
	CurveCertificate     *czmq.Cert
	rpc                  *rpc.Server
	worker               *Worker
}

func NewJSONRPCWorker(broker_endpoint, service string) (j *JSONRPCWorker) {
	j = &JSONRPCWorker{}
	j.rpc = rpc.NewServer()
	j.worker = NewWorker(broker_endpoint, service)
	return j
}

func (j *JSONRPCWorker) Register(handler interface{}) error {
	err := j.rpc.Register(handler)
	if err != nil {
		return err
	}

	typ := reflect.TypeOf(handler)
	rcvr := reflect.ValueOf(handler)
	x := reflect.Indirect(rcvr).Type().Name()
	log.Printf("Register: %s / %s / %s", typ, rcvr, x)

	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type
		log.Printf("Method: %s -- %s", mtype, method.Name)
		log.Printf("Arg: %s", mtype.In(1))
		log.Printf("Reply: %s", mtype.In(2))
	}

	return nil
}

func (j *JSONRPCWorker) Run() {
	if len(j.CurveServerPublicKey) > 0 {
		j.worker.CurveServerPublicKey = j.CurveServerPublicKey
	}
	if j.CurveCertificate != nil {
		j.worker.CurveCertificate = j.CurveCertificate
	}
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
