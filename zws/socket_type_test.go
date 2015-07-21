package zws

import (
	"testing"
)

var VALID_TYPES = []SocketType{REQ}

func TestSocketTypeisValidTrue(t *testing.T) {
	for _, socket_type := range VALID_TYPES {
		if !socket_type.isValid() {
			t.Errorf("Socket type %d[%s] should be valid, but isValid() returned false")
			return
		}
	}
}

func TestSocketTypeisValidFalse(t *testing.T) {
	types := []SocketType{0, 100}

	for _, socket_type := range types {
		if socket_type.isValid() {
			t.Errorf("Socket type %d[%s] should be invalid valid, but isValid() returned true")
			return
		}
	}
}

func TestSocketTypeCreateReturnsSock(t *testing.T) {
	endpoint := "inproc://example"
	for _, socket_type := range VALID_TYPES {
		s, err := socket_type.Create(endpoint)
		if err != nil {
			t.Errorf("%s.Create(`%s`) failed: %s", socket_type, endpoint, err)
			return
		}
		if s == nil {
			t.Errorf("%s.Create(`%s`) returnd nil socket", socket_type, endpoint)
			return
		}
	}
}
