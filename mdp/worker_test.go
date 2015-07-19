package mdp

import (
	"fmt"
	czmq "github.com/zeromq/goczmq"
	"testing"
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

	// Skip the first frame since that's the DEALER/ROUTER id
	err = validateWorkerReady(frames[1:], service)
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

	// Skip the first frame since that's the DEALER/ROUTER id
	err = validateWorkerReady(frames[1:], service)
	if err != nil {
		t.Errorf("Error reading READY message: %s\n", err)
		return
	}

	frames, err = sock.RecvMessage()
	err = validateWorkerHeartbeat(frames[1:])
	if err != nil {
		t.Errorf("Error in HEARTBEAT: %s\n", err)
		return
	}
}

type HelloGreeter struct{}

func (h *HelloGreeter) Request(request [][]byte) (response [][]byte, err error) {
	response = append(response, []byte("Nice to meet you!"))
	return
}

func (h *HelloGreeter) Heartbeat()  {}
func (h *HelloGreeter) Disconnect() {}

func TestWorkerRun(t *testing.T) {
	broker := fmt.Sprintf("inproc://%s", randomHex())
	service := randomHex()
	client := randomHex()

	sock, err := czmq.NewRouter(broker)
	if err != nil {
		t.Errorf("NewRouter(%v) failed", broker)
		return
	}

	go func(broker, service string) {
		NewWorker(broker, service).Run(&HelloGreeter{})
	}(broker, service)

	frames, _ := sock.RecvMessage()
	err = validateWorkerReady(frames[1:], service)
	if err != nil {
		t.Errorf("Error reading READY message: %s\n", err)
		return
	}

	err = sock.SendMessage([][]byte{
		frames[0], // router/dealer ID
		[]byte{},
		MDP_WORKER,
		[]byte{byte(C_REQUEST)},
		[]byte(client),
		[]byte{}, // SPEC Frame 4: Empty (zero bytes, envelope delimiter)
		[]byte("hello world"),
	})
	reply, err := sock.RecvMessage()
	if err != nil {
		t.Errorf("Error receiving reply: %s\n", err)
	}
	err = validateWorkerReply(reply[1:], []byte(client)) // [1:] to skip the router/dealer id
	if err != nil {
		t.Errorf("Reply validation failed: %s\n", err)
	}
}
