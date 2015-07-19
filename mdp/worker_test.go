package mdp

import (
	czmq "github.com/zeromq/goczmq"
  "testing"
  "bytes"
  "fmt"
)

func validateReadyMessage(t *testing.T, frames [][]byte, service string) error {
  // Since this is a ROUTER socket, the first frame is the client id, so the MDP frames start at frames[1]
  //for i, x := range frames {
    //fmt.Printf("%d: %v\n", i, string(x))
  //}
  if count := len(frames); count < 5 {
    return fmt.Errorf("Expected at least 5 frames, got %d", count)
  }

  if length := len(frames[1]); length != 0 {
    return fmt.Errorf("First frame must be empty")
  }

  if !bytes.Equal(frames[2], MDP_WORKER) {
    return fmt.Errorf("Second frame must be %s, got %s", string(MDP_WORKER), string(frames[2]))
  }

  if !bytes.Equal(frames[3], []byte{C_READY}) {
    return fmt.Errorf("Third frame must be 0x%02x (READY command).", C_READY)
  }

  if !bytes.Equal(frames[4], []byte(service)) {
    return fmt.Errorf("Service mismatch. Expected %s but received %s", service, string(frames[4]))
  }
  return nil
}


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
    fmt.Errorf("NewRouter(%v) failed", broker)
    return
  }

  go func(broker, service string) {
    NewWorker(broker, service).ensure_connected()
  }(broker, service)

  // Worker should send READY as soon as it connects

  // SPEC: Frame 0: Empty frame
  // SPEC: Frame 1: "MDPW01" (six bytes, representing MDP/Worker v0.1)
  // SPEC: Frame 2: 0x01 (one byte, representing READY)
  // SPEC: Frame 3: Service name (printable string)
  frames, err := sock.RecvMessage()
  if err != nil {
    t.Errorf("%s", err)
    return
  }

  err = validateReadyMessage(t, frames, service)
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
    fmt.Errorf("NewRouter(%v) failed", broker)
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
    fmt.Errorf("Error reading READY message: %s\n", err)
    return
  }
  err = validateReadyMessage(t, frames, service)
  if err != nil {
    fmt.Errorf("Error reading READY message: %s\n", err)
    return
  }

  // SPEC: Frame 0: Empty frame
  // SPEC: Frame 1: "MDPW01" (six bytes, representing MDP/Worker v0.1)
  // SPEC: Frame 2: 0x04 (one byte, representing HEARTBEAT)
  frames, err = sock.RecvMessage()

  // Since this is a ROUTER socket, the first frame is the client id, so the MDP frames start at frames[1]

  //for i, x := range frames {
    //fmt.Printf("%d: %v\n", i, string(x))
  //}
  if count := len(frames); count != 4 {
    t.Errorf("Expected exactly 4 frames, got %d", count)
    return
  }

  if length := len(frames[1]); length != 0 {
    t.Errorf("First frame must be empty")
    return
  }

  if !bytes.Equal(frames[2], MDP_WORKER) {
    t.Errorf("Second frame must be %s, got %s", string(MDP_WORKER), string(frames[2]))
    return
  }

  if !bytes.Equal(frames[3], []byte{C_HEARTBEAT}) {
    t.Errorf("Third frame must be 0x%02x (HEARTBEAT command).", C_HEARTBEAT)
    return
  }
}
