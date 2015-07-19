package mdp

import (
	"testing"
)

func TestConstants(t *testing.T) {
	if string(MDP_CLIENT) != "MDPC01" {
		t.Errorf("MDP_CLIENT is not correct")
	}

	if string(MDP_WORKER) != "MDPW01" {
		t.Errorf("MDP_WORKER is not correct")
	}

	if C_READY != 0x01 {
		t.Errorf("READY command has wrong value")
	}
	if C_REQUEST != 0x02 {
		t.Errorf("REQUEST command has wrong value")
	}
	if C_REPLY != 0x03 {
		t.Errorf("REPLY command has wrong value")
	}
	if C_HEARTBEAT != 0x04 {
		t.Errorf("HEARTBEAT command has wrong value")
	}

	if C_DISCONNECT != 0x05 {
		t.Errorf("DISCONNECT command has wrong value")
	}
}
