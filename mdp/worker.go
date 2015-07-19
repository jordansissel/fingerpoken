package mdp

import (
	"fmt"
	czmq "github.com/jordansissel/goczmq"
	"log"
	"time"
)

// TODO(sissel): turn this into an interface?
type Worker struct {
	sock               *czmq.Sock
	broker             string
	service            string
	HeartbeatFrequency time.Duration
}

func NewWorker(broker string, service string) (w *Worker) {
	w = &Worker{broker: broker, service: service}
	w.HeartbeatFrequency = 5 * time.Second
	return
}

type RequestHandler interface {
	Request([][]byte) ([][]byte, error)
	Heartbeat()
	Disconnect()
}

func (w *Worker) Run(requestHandler RequestHandler) error {
	w.ensure_connected()

	for {
		// TODO(sissel): poll and send heartbeats
		client, command, body, err := w.readRequest()
		if err != nil {
			log.Printf("Worker: Error reading request: %s\n", err)
			w.Reset()
			continue
		}
		switch command {
		case C_HEARTBEAT:
			requestHandler.Heartbeat()
		case C_DISCONNECT:
			requestHandler.Disconnect()
			w.sock.Destroy()
		case C_REQUEST:
			// The spec supports multiple frames for a request message. Let's support that.
			reply_body, err := requestHandler.Request(body)
			if err != nil {
				log.Printf("Worker: Error handling request: %s\n", err)
				w.Reset()
				continue
			}

			err = w.sock.SendMessage(append([][]byte{
				[]byte{},
				MDP_WORKER,
				[]byte{byte(C_REPLY)},
				client,
				[]byte{}, // SPEC Frame 4: Empty (zero bytes, envelope delimiter)
			}, reply_body...),
			)
			if err != nil {
				log.Printf("Worker: Error sending reply: %s\n", err)
				w.Reset()
				continue
			}

		default:
			log.Printf("Worker: Got an invalid command from broker. Will reset connection. (command: %v)", command)
			w.Reset()
		}
	}
}

func (w *Worker) readRequest() (client []byte, command Command, body [][]byte, err error) {
	frames, err := w.sock.RecvMessage()
	if err != nil {
		return
	}
	for i, x := range frames {
		log.Printf("Worker(via Broker): frame %d: %v (%s)\n", i, x, string(x))
	}

	err = validateWorkerRequest(frames[:])
	if err != nil {
		err = fmt.Errorf("Got an invalid worker request in request. Will reset connection. %s", err)
		return
	}

	command = Command(frames[2][0])
	client = frames[3]
	body = frames[5:]
	return
}

func (w *Worker) Reset() {
	w.sock.Destroy()
	w.sock = nil
	w.ensure_connected()
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
		[]byte{},              // SPEC: Frame 0: Empty frame
		MDP_WORKER,            // SPEC: Frame 1: "MDPW01" (six bytes, representing MDP/Worker v0.1)
		[]byte{byte(C_READY)}, // SPEC: Frame 2: 0x01 (one byte, representing READY)
		[]byte(w.service),     // SPEC: Frame 3: Service name (printable string)
	}

	err = w.sock.SendMessage(message[:])
	// SPEC: There is no response to a READY. The worker SHOULD assume the registration succeeded.
	if err != nil {
		return err
	}
	return nil
}

func (w *Worker) sendHeartbeat() (err error) {
	err = w.sock.SendMessage(M_HEARTBEAT[:])
	if err != nil {
		return err
	}
	return nil
}
