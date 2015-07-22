package mdp

import (
	"time"
)

type workerEntry struct {
	expiration            time.Time
	service               string
	address               []byte
	nextSendHeartbeatTime time.Time
}

func (entry *workerEntry) recordHeartbeat(expiration time.Time) {
	entry.expiration = expiration
}
