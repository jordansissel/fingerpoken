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
