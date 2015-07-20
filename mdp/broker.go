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

	HeartbeatCallback func(*WorkerEntry)
}

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
		s := czmqPollerSafeWait(poller, durationInMilliseconds(b.HeartbeatInterval))
		if s != nil {
			b.once(s)
		}

		for _, entry := range b.workers {
			if time.Now().After(entry.expiration) {
				// We haven't seen this worker for too long, let's delete it.
				log.Printf("Worker hasn't been seen in a while, removing it.")
				b.removeWorker(entry.address)
			} else {
				// Healthy worker. Let's tell the worker we're still here if needed.
				if time.Now().After(entry.nextSendHeartbeatTime) {
					b.sendHeartbeat(entry)
				}
			}
		}
	}
}

func (b *Broker) once(sock *czmq.Sock) {
	frames, _ := sock.RecvMessage()

	// Min. 4 frames   ||  Frame 1 is empty   || Frame 2 is MDPW01   || Frame 3 client addr || Frame 4 is empty
	if len(frames) < 4 || len(frames[1]) != 0 || len(frames[2]) != 6 {
		// SPEC: broker SHOULD respond to invalid messages by dropping them and treating that peer as invalid.
		log.Printf("Broker: Got an invalid message, skipping this message.")
		log.Printf("Frames(%d), frame1(%d), frame2(%d), frame3(%d), frame4(%d)\n", len(frames), len(frames[1]), len(frames[2]), len(frames[3]), len(frames[4]))
		return
	}

	address := frames[0]
	if bytes.Equal(frames[2], mdp_WORKER) {
		b.handleWorker(address, frames[1:])
	} else {
		b.handleClient(address, frames[1:])
	}
}

func (b *Broker) handleWorker(address []byte, frames [][]byte) {
	//for i, x := range frames { log.Printf("Broker(via Worker): frame %d: %v (%s)\n", i, x, string(x)) }
	if len(frames) < 3 || len(frames[1]) != 6 || len(frames[2]) != 1 {
		log.Printf("Broker: Got an invalid worker message. Dropping.")
		return
	}
	// precondition: this is a well formed message that has enough frames and of the right size.
	cmd := command(frames[2][0])

	// SPEC The broker MUST respond to any valid but unexpected cmd by sending DISCONNECT
	//
	// We'll use the presence check `known` in the workers map to determine if the worker has
	// sent in a READY cmd. If it hasn't, then it's not a valid worker and we should
	// tell it to disconnect and ignore it.
	entry, known := b.getWorkerByAddress(address)
	if !known {
		if cmd == c_READY { // New worker. First cmd must be READY.
			b.addWorker(address, string(frames[3]))
		} else {
			log.Printf("Broker: Received cmd from unknown worker (one that has not sent READY yet). Will disconnect it.")
			// TODO(sissel): Send disconnect
		}
		return
	}

	entry.recordHeartbeat(b.nextExpiration())

	switch cmd {
	case c_HEARTBEAT: // Nothing to do
		if b.HeartbeatCallback != nil {
			b.HeartbeatCallback(entry)
		}
	case c_DISCONNECT:
		log.Printf("Broker: Received disconnect from worker %v", address)
		b.removeWorker(address)
	case c_REQUEST:
		log.Printf("Broker: received REQUEST from a worker. This is not valid. Workers don't make requests.")
		b.removeWorker(address)
		// TODO(sissel): Send Disconnect
		return
	case c_REPLY:
		destination := frames[3]
		log.Printf("Broker: Received reply from worker. Destination %v", destination)
		replyheader := [4][]byte{
			destination,
			[]byte{},
			mdp_CLIENT,
			[]byte(entry.service),
		}
		client_reply := append(replyheader[:], frames[4:]...)
		//for i, x := range client_reply { log.Printf("Broker(reply to Client): frame %d: %v (%s)\n", i, x, string(x)) }
		err := b.sock.SendMessage(client_reply)
		if err != nil {
			log.Printf("Broker: Error forwarding reply to client: %s\n", err)
		}
	}
} // handleWorker

func (b *Broker) addWorker(address []byte, service string) {
	log.Printf("New worker providing `%s`\n", service)
	key := base64.StdEncoding.EncodeToString(address)
	b.workers[key] = &WorkerEntry{
		expiration: b.nextExpiration(),
		service:    service, // SPEC: Frame 3: Service name (printable string)
		address:    address,
	}
}

func (b *Broker) getWorkerByAddress(address []byte) (*WorkerEntry, bool) {
	key := base64.StdEncoding.EncodeToString(address)
	entry, ok := b.workers[key]
	return entry, ok
}

func (b *Broker) getWorkerByService(service string) (*WorkerEntry, bool) {
	for _, e := range b.workers {
		if e.service == service {
			return e, true
		}
	}
	return nil, false
}

func (b *Broker) removeWorker(address []byte) {
	key := base64.StdEncoding.EncodeToString(address)
	delete(b.workers, key)
}

func (b *Broker) handleClient(address []byte, frames [][]byte) {
	//for i, x := range frames { log.Printf("Broker(via Client %v): frame %d: %v (%s)\n", address, i, x, string(x)) }
	err := validateClientHeader(frames)
	if err != nil {
		log.Printf("Broker: Received invalid request from client. Dropping message.\n")
		return
	}

	cmd := command(frames[2][0])
	if cmd != c_REQUEST {
		log.Printf("Broker: Received invalid cmd (%s) from client. Dropping message.\n", cmd)
		return
	}

	service := string(frames[3])

	// TODO(sissel): Find a worker providing the given service
	entry, found := b.getWorkerByService(service)

	if !found {
		// SPEC: The broker SHOULD queue client requests for which no service has
		// SPEC: been registered and SHOULD expire these requests after a
		// SPEC: reasonable and configurable time if no service has been registered.
		// TODO(sissel): Implement queuing?
		log.Printf("Broker: Received request for service `%s` which has no workers. Dropping request.", service)
		return
	}

	// TODO(sissel): Send the worker the request
	worker_header := [][]byte{
		entry.address,
		[]byte{},
		mdp_WORKER,
		[]byte{byte(c_REQUEST)},
		address,
		[]byte{},
	}
	message := append(worker_header[:], frames[3:]...)
	//for i, x := range message { log.Printf("Broker(-> Worker): frame %d: %v (%s)\n", i, x, string(x)) }
	err = b.sock.SendMessage(message)
	if err != nil {
		log.Printf("Broker: Error forwarding client request to a worker: %s\n", err)
	}
}

func (b *Broker) nextExpiration() time.Time {
	// SPEC: A peer MUST consider the other peer "disconnected" if no heartbeat arrives within some multiple of that interval (usually 3-5).
	return time.Now().Add(time.Duration(int64(b.HeartbeatInterval) * b.MaxMissedHeartbeats))
}

func (b *Broker) sendHeartbeat(entry *WorkerEntry) {
	heartbeat := [][]byte{entry.address}
	heartbeat = append(heartbeat[:], m_HEARTBEAT[:]...)
	err := b.sock.SendMessage(heartbeat)
	if err != nil {
		log.Printf("Broker: Error sending heartbeat", err)
		// TODO(sissel): what should we do?
	}
}

func (b *Broker) sendDisconnect(entry *WorkerEntry) {
	disconnect := [][]byte{entry.address}
	disconnect = append(disconnect[:], m_DISCONNECT[:]...)
	err := b.sock.SendMessage(disconnect)
	if err != nil {
		log.Printf("Broker: Error sending disconnect", err)
		// TODO(sissel): what should we do?
	}
}
