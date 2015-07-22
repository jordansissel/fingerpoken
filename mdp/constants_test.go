package mdp

import (
	"testing"
)

func TestConstants(t *testing.T) {
	if string(mdp_CLIENT) != "MDPC01" {
		t.Errorf("mdp_CLIENT is not correct")
	}

	if string(mdp_WORKER) != "MDPW01" {
		t.Errorf("mdp_WORKER is not correct")
	}

	if c_READY != 0x01 {
		t.Errorf("READY cmd has wrong value")
	}
	if c_REQUEST != 0x02 {
		t.Errorf("REQUEST cmd has wrong value")
	}
	if c_REPLY != 0x03 {
		t.Errorf("REPLY cmd has wrong value")
	}
	if c_HEARTBEAT != 0x04 {
		t.Errorf("HEARTBEAT cmd has wrong value")
	}

	if c_DISCONNECT != 0x05 {
		t.Errorf("DISCONNECT cmd has wrong value")
	}
}
