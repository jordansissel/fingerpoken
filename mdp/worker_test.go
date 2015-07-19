package mdp

import (
	czmq "github.com/zeromq/goczmq"
  "testing"
  "fmt"
)

func TestWorkerCreation(t *testing.T) {
  broker := fmt.Sprintf("inproc://%s", randomHex())
  service := randomHex()
  w := NewWorker(broker, service)
  err := w.ensure_connected()
  if err != nil {
    t.Errorf("Worker#ensure_connected failed: %s\n", err)
    return
  }
}

func TestWorkerReadyMessage(t *testing.T) {
  broker := fmt.Sprintf("inproc://%s", randomHex())
  service := randomHex()

  sock, err := czmq.NewRouter(broker)
  if err != nil {
    t.Errorf("NewRouter(%v) failed", broker)
    return
  }

  go func(broker, service string) {
    NewWorker(broker, service).ensure_connected()
  }(broker, service)

  // Worker should send READY as soon as it connects
  frames, err := sock.RecvMessage()
  if err != nil {
    t.Errorf("%s", err)
    return
  }

  err = validateReadyMessage(frames, service)
  if err != nil {
    t.Errorf("%s", err)
    return
  }
}

func TestWorkerHeartbeatMessage(t *testing.T) {
  broker := fmt.Sprintf("inproc://%s", randomHex())
  service := randomHex()

  sock, err := czmq.NewRouter(broker)
  if err != nil {
    t.Errorf("NewRouter(%v) failed", broker)
    return
  }

  go func(broker, service string) {
    w := NewWorker(broker, service)
    w.ensure_connected()
    w.sendHeartbeat()
  }(broker, service)

  // Worker should send READY as soon as it connects
  frames, err := sock.RecvMessage()
  if err != nil {
    t.Errorf("Error reading READY message: %s\n", err)
    return
  }
  err = validateReadyMessage(frames, service)
  if err != nil {
    t.Errorf("Error reading READY message: %s\n", err)
    return
  }

  frames, err = sock.RecvMessage()
  err = validateHeartbeat(frames)
  if err != nil {
    t.Errorf("Error in HEARTBEAT: %s\n", err)
    return
  }
}
