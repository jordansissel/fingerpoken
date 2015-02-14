package main

import (
	//"golang.org/x/net/websocket"
	"log"
  "os/exec"
	"net/http"
)

func main() {
	// Simple static webserver:
  mux := http.NewServeMux()

  mux.HandleFunc("/hello", func(response http.ResponseWriter, request *http.Request) {
    response.Write([]byte("Hello world"))
  })

  mux.HandleFunc("/lirc/once/onkyo/GAME", func(response http.ResponseWriter, request *http.Request) {
    //irsend SEND_ONCE onkyo "BD/DVD
    cmd := exec.Command("irsend", "SEND_ONCe", "onkyo", "GAME")
    err := cmd.Run()
    if err != nil {
      log.Printf("Failed to run command: %v", err)
      http.Error(response, "Something went wrong", 503)
      return
    }

    response.Write([]byte("OK"))
  })

  mux.Handle("/", http.FileServer(http.Dir("./app")))
	log.Fatal(http.ListenAndServe(":8080", mux))
}
