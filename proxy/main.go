package main

import (
  "github.com/gorilla/mux"
  "github.com/gorilla/websocket"
  "net/http"
  "fmt"
  consul "github.com/hashicorp/consul/api"
  czmq "github.com/zeromq/goczmq"
)

// Many gorilla websocket examples have a global 'upgrader' but
// XXX: I feel like this isn't threadsafe...
var upgrader = websocket.Upgrader{
  ReadBufferSize:  4096,
  WriteBufferSize: 4096,
}

func serveWebSocket(w http.ResponseWriter, r *http.Request) {
  conn, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    fmt.Println(err)
    return
  }

  client, _ := consul.NewClient(consul.DefaultConfig())
  services, _, _ := client.Catalog().Service("rpc", "", nil)

  var endpoints []*czmq.Sock;
  for _, service := range services {
    endpoint := fmt.Sprintf("tcp://%s:%d", service.Address, service.ServicePort)
    fmt.Println(endpoint)
    req, _ := czmq.NewReq(endpoint)
    endpoints = append(endpoints, req)
  }

  err = nil
  for err == nil {
    mtype, payload, err := conn.ReadMessage()
    fmt.Printf("%s Payload [%v]: %s\n", mtype, err, string(payload))
    for _, endpoint := range endpoints {
      err = endpoint.SendFrame(payload[:], 0)
      if err != nil {
        fmt.Printf("endpoint.SendMessage fail: %s\n", err)
        panic("!")
      }
      response, err := endpoint.RecvMessage()
      if err != nil {
        fmt.Printf("endpoint.RecvMessage fail: %s\n", err)
        panic("!")
      }
      fmt.Printf("Resp: %v\n", string(response[0]))
    }
  }
}

func main() {
  r := mux.NewRouter()
  r.Handle("/", http.FileServer(http.Dir("./static")))
  r.PathPrefix("/js/").Handler(http.StripPrefix("/js", http.FileServer(http.Dir("./static/js/"))))
  r.HandleFunc("/ws", serveWebSocket)
  http.Handle("/", r)
  http.ListenAndServe(":8000", nil)
}
