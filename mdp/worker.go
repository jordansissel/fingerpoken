package mdp

import (
	czmq "github.com/zeromq/goczmq"
  "time"
)

// TODO(sissel): turn this into an interface?
type Worker struct {
  sock *czmq.Sock
  broker string
  service string
  HeartbeatFrequency time.Duration
}

func NewWorker(broker string, service string) (w *Worker){
  w = &Worker{broker: broker, service: service}
  w.HeartbeatFrequency = 5 * time.Second
  return
}

func (w *Worker) ensure_connected() error {
	if w.sock != nil {
		return nil
	}

	var err error
	w.sock, err = czmq.NewDealer(w.broker)
	if err != nil {
		return err
	}

  err = w.sendReady()
  if err != nil {
    return err
  }
	return nil
}

func (w *Worker) sendReady() (err error) {
  message := [4][]byte{
    []byte{}, // SPEC: Frame 0: Empty frame
    MDP_WORKER, // SPEC: Frame 1: "MDPW01" (six bytes, representing MDP/Worker v0.1)
    []byte{C_READY}, // SPEC: Frame 2: 0x01 (one byte, representing READY)
    []byte(w.service), // SPEC: Frame 3: Service name (printable string)
  }

  err = w.sock.SendMessage(message[:])
  // SPEC: There is no response to a READY. The worker SHOULD assume the registration succeeded.
  if err != nil {
    return err
  }
  return nil
}

func (w *Worker) sendHeartbeat() (err error) {
  message := [3][]byte{
    []byte{}, // SPEC: Frame 0: Empty frame
    MDP_WORKER, // SPEC: Frame 1: "MDPW01" (six bytes, representing MDP/Worker v0.1)
    []byte{C_HEARTBEAT}, // SPEC: Frame 2: 0x04 (one byte, representing HEARTBEAT)
  }

  err = w.sock.SendMessage(message[:])
  if err != nil {
    return err
  }
  return nil
}
