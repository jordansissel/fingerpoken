package util

import (
  "net/rpc"
  "net/rpc/jsonrpc"
  "bytes"
  "fmt"
  "io/ioutil"
  consul "github.com/hashicorp/consul/api"
  czmq "github.com/zeromq/goczmq"
)

type ZJServer struct {
  repsock *czmq.Sock
  endpoint string
  rpc *rpc.Server
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
  err = registerConsulService(client, "rpc", zj.endpoint)
  return
}

func (zj *ZJServer) Register(handler interface{}) error {
  // TODO(sissel): record the methods supported by `handler` and expose that as a ListMethods call
  return zj.rpc.Register(handler)
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

  input := bytes.NewBuffer(message[0])
  output := &bytes.Buffer{}
  codec := jsonrpc.NewServerCodec(buffer { Reader: input, Writer: output, Closer: ioutil.NopCloser(nil) })
  err = zj.rpc.ServeRequest(codec)
  if err != nil {
    fmt.Printf("ServeRequest err: %s\n", err);
  }

  response := [1][]byte{ output.Bytes() }
  zj.repsock.SendMessage(response[:])
  return
}
