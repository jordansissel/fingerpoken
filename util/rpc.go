package util

import (
	"bytes"
	consul "github.com/hashicorp/consul/api"
	//czmq "github.com/zeromq/goczmq"
	czmq "github.com/zeromq/goczmq"
	"io/ioutil"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type ZJServer struct {
	repsock  *czmq.Sock
	endpoint string
	rpc      *rpc.Server
}

func NewZJServer() (*ZJServer, error) {
	zjs := &ZJServer{}
	err := zjs.Bind("tcp://*:*")
	if err != nil {
		return nil, err
	}
	zjs.rpc = rpc.NewServer()
	return zjs, nil
}

func (zj *ZJServer) Bind(endpoint string) (err error) {
	zj.repsock, err = czmq.NewRep(endpoint)
	if err != nil {
		return
	}

	zj.endpoint = zj.repsock.LastEndpoint()
	return
}

func (zj *ZJServer) RegisterWithConsul(client *consul.Client) (err error) {
	err = ConsulRegisterService(client, "rpc", zj.endpoint)
	return
}

func (zj *ZJServer) Register(handler interface{}) (err error) {
	// TODO(sissel): record the methods supported by `handler` and expose that as a ListMethods call
	err = zj.rpc.Register(handler)
	return
}

func (zj *ZJServer) Loop() (err error) {
	for err == nil {
		err = zj.Once()
	}
	return
}

func (zj *ZJServer) Once() (err error) {
	message, err := zj.repsock.RecvMessage()
	if err != nil {
		return
	}
	//fmt.Printf("> %s\n", string(message[0]))

	input := bytes.NewBuffer(message[0])
	output := &bytes.Buffer{}
	codec := jsonrpc.NewServerCodec(Buffer{Reader: input, Writer: output, Closer: ioutil.NopCloser(nil)})
	err = zj.rpc.ServeRequest(codec)
	if err != nil {
		// If the rpc call fails, do we want to return an error? I don't think so.
		//fmt.Printf("ServeRequest err: %s\n", err);
	}

	zj.repsock.SendFrame(output.Bytes(), 0)
	return
}
