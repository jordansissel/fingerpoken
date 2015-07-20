package mdp

import (
	"time"
)

type WorkerEntry struct {
	expiration            time.Time
	service               string
	address               []byte
	nextSendHeartbeatTime time.Time
}

func (entry *WorkerEntry) recordHeartbeat(expiration time.Time) {
	entry.expiration = expiration
}
