package mdp

import (
	"bytes"
	"encoding/base64"
	czmq "github.com/zeromq/goczmq"
	"log"
	"time"
)

type Broker struct {
	sock                *czmq.Sock
	endpoint            string
	workers             map[string]*WorkerEntry
	HeartbeatInterval   time.Duration
	MaxMissedHeartbeats int64
}

type WorkerEntry struct {
	expiration time.Time
	service    string
}

type Address []byte

func NewBroker(endpoint string) (b *Broker, err error) {
	b = &Broker{
		endpoint:            endpoint,
		HeartbeatInterval:   5 * time.Second,
		MaxMissedHeartbeats: 3,
	}
	b.workers = make(map[string]*WorkerEntry)

	b.sock, err = czmq.NewRouter(endpoint)
	if err != nil {
		return
	}

	// Anything else to do?
	return
}

func (b *Broker) Run() {
	poller, err := czmq.NewPoller(b.sock)
	if err != nil {
		log.Fatalf("czmq.NewPoller failed: %s\n", err)
	}

	for {
		sock := poller.Wait(1000)
		if sock == nil {
			// Timeout, do any maintenance tasks?
			// TODO(sissel): Purge any expired workers
		} else {
			b.Once(sock)
		}
	}
}

func (b *Broker) Once(sock *czmq.Sock) {
	frames, _ := sock.RecvMessage()
	for i, x := range frames {
		log.Printf("broker %d: %v (%s)\n", i, x, string(x))
	}

	if len(frames) < 4 || len(frames[1]) != 0 || len(frames[2]) != 6 || len(frames[3]) != 1 {
		// SPEC: broker SHOULD respond to invalid messages by dropping them and treating that peer as invalid.
		log.Printf("Got an invalid message, skipping this message.")
		return
	}

	address := frames[0]
	if bytes.Equal(frames[2], MDP_WORKER) {
		b.handleWorker(address, frames[1:])
	} else {
		b.handleClient(address, frames[1:])
	}
}

func (b *Broker) handleWorker(address []byte, frames [][]byte) {
	// precondition: this is a well formed message that has enough frames and of the right size.
	worker_key := base64.StdEncoding.EncodeToString(address)
	entry, known := b.workers[worker_key]

	// SPEC The broker MUST respond to any valid but unexpected command by sending DISCONNECT
	// We'll use the presence check `known` in the workers map to determine if the worker has
	// sent in a READY command. If it hasn't, then it's not a valid worker and we should
	// tell it to disconnect and ignore it.
	command := Command(frames[2][0])
	switch command {
	case C_READY:
		b.workers[worker_key] = &WorkerEntry{
			expiration: b.nextExpiration(),
			service:    string(frames[3]), // SPEC: Frame 3: Service name (printable string)
		}
		return
	default:
		if !known {
			log.Printf("Received command from unknown worker (one that has not sent READY yet). Will disconnect it.")
			// TODO(sissel): Send disconnect
			return
		}
	}

	entry.heartbeat(b.nextExpiration())

	switch command {
	case C_HEARTBEAT: // Nothing to do
	case C_DISCONNECT:
		log.Printf("Received disconnect from worker %s", worker_key)
		delete(b.workers, worker_key)
	case C_REQUEST:
		log.Printf("Broker received REQUEST from a worker. This is not valid. Workers don't make requests.")
		// TODO(sissel): Send Disconnect
		delete(b.workers, worker_key)
		return
	case C_REPLY:
		log.Printf("Received reply from worker")
	}
}

func (b *Broker) handleClient(address []byte, frames [][]byte) {
	panic("IMPLEMENT ME")
}

func (b *Broker) nextExpiration() time.Time {
	// SPEC: A peer MUST consider the other peer "disconnected" if no heartbeat arrives within some multiple of that interval (usually 3-5).
	return time.Now().Add(time.Duration(int64(b.HeartbeatInterval) * b.MaxMissedHeartbeats))
}

func (entry *WorkerEntry) heartbeat(expiration time.Time) {
	entry.expiration = expiration
}
