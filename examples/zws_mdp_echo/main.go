package main

import (
  "github.com/jordansissel/fingerpoken/zws"
  "log"
  "net/http"
  "github.com/jordansissel/fingerpoken/mdp"
)

func main() {
  broker_endpoint := "inproc://fancy-req"
  service := "echo"

  b, err := mdp.NewBroker(broker_endpoint)
  if err != nil {
    log.Fatalf("NewBroker(%s) failed: %s", broker_endpoint, err)
  }
  port, err := b.Bind("tcp://*:*")
  log.Printf("Broker available at binding: tcp://x:%d", port)
  go b.Run()

  // Start a worker.
	w := mdp.NewWorker(broker_endpoint, service)
	go w.Run(&echoWorker{})

  RunHTTP(":8111", "/zws/1.0")
}

func RunWorker(endpoint string) {
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
