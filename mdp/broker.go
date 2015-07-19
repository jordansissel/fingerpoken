package mdp

import (
	czmq "github.com/zeromq/goczmq"
)

type Broker struct {
  sock *czmq.Sock
  endpoint string
}

func NewBroker(endpoint string) (b *Broker) {
  b = &Broker{endpoint: endpoint}
  return
}



