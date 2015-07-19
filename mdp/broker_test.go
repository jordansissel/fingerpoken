package mdp

import (
	"fmt"
	"testing"
)

func TestBrokerWorkerIntegration(t *testing.T) {
	broker_address := fmt.Sprintf("inproc://%s", randomHex())
	service := randomHex()
	w := NewWorker(broker_address, service)
	b, _ := NewBroker(broker_address)
	go func() { b.Run() }()
	err := w.ensure_connected()
	if err != nil {
		t.Errorf("Worker#ensure_connected failed: %s\n", err)
		return
	}
}
