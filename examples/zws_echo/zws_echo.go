package main

import (
  "github.com/jordansissel/fingerpoken/zws"
  "log"
	czmq "github.com/zeromq/goczmq"
  "net/http"
)

func main() {
  endpoint := "inproc://fancy-req"
  go RunRep(endpoint)
  RunHTTP(":8111", "/zws/1.0")
}

func RunRep(endpoint string) {
  rep, err := czmq.NewRep(endpoint)
  if err != nil {
    panic(err)
  }

  for {
    message, flags, err := rep.RecvFrame()
    if err != nil {
      log.Fatalf("rep.RecvMessage failed: %s", err)
    }

    // Echo it back.
    err = rep.SendFrame(message, flags)
    if err != nil {
      log.Fatalf("rep.SendMessage failed: %s", err)
    }
  }
}

func RunHTTP(address, path string) {
  mux := http.NewServeMux()
  mux.HandleFunc(path, zws.HandleZWS)
  log.Printf("Listening on http://%s%s", address, path)
  server := http.Server{Addr: address, Handler: mux}
  err := server.ListenAndServe()
  if err != nil {
    log.Fatalf("server.ListenAndServe failed: %s", err)
  }
}
