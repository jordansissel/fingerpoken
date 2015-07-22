// This file is part of fingerpoken
// Copyright (C) 2015 Jordan Sissel
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
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
