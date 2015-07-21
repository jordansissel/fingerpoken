package zws

import (
	"github.com/gorilla/websocket"
	//"github.com/gorilla/mux"
	czmq "github.com/zeromq/goczmq"
  "net/http"
  "fmt"
)

type ZWS struct {
  ws *websocket.Conn
  zmq *czmq.Sock
  socket_type SocketType
}

func NewZWS(endpoint_prefix string, w http.ResponseWriter, r *http.Request) (*ZWS, error) {
  query := r.URL.Query()
  socket_type_name, ok := query["type"]
  if !ok {
    return nil, fmt.Errorf("Request is missing `?type=SOCKET_TYPE` from path.")
  }

  number, err := ParseSocketType(socket_type_name[0])
  if err != nil {
    return nil, fmt.Errorf("`type=%s`: %s", socket_type_name[0], err)
  }
  socket_type := SocketType(number)

  if !socket_type.isValid() {
    return nil, fmt.Errorf("`type=%s` is not valid (%s)", socket_type_name, socket_type)
  }

  zmq, err := socket_type.Create(fmt.Sprintf("%s-%s", endpoint_prefix, socket_type.EndpointSuffix()))
  if err != nil {
    return nil, fmt.Errorf("socket creation failed: %s", err)
  }

  ws, err := upgrade(w, r)
  if err != nil {
    return nil, err
  }

  zws := &ZWS{
    ws: ws,
    socket_type: socket_type,
    zmq: zmq,
  }
  return zws, nil
}

