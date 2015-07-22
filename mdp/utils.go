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
	"bytes"
	"fmt"
	czmq "github.com/zeromq/goczmq"
	"math/rand"
	"time"
)

func randomHex() (value string) {
	var length = rand.Int31n(10) + 5
	for x := int32(0); x < length; x += 1 {
		value = fmt.Sprintf("%s%02x", value, rand.Int31n(256))
	}
	return
}

func validateClientHeader(frames [][]byte) error {
	if length := len(frames[0]); length != 0 {
		return fmt.Errorf("First frame must be empty")
	}

	if !bytes.Equal(frames[1], mdp_CLIENT) {
		return fmt.Errorf("Second frame must be %s, got %s", string(mdp_CLIENT), string(frames[2]))
	}
	return nil
}

func validateClientRequest(frames [][]byte, service string) error {
	if count := len(frames); count < 4 {
		return fmt.Errorf("Expected at least 4 frames, got %d", count)
	}

	if err := validateClientHeader(frames); err != nil {
		return err
	}

	if !bytes.Equal(frames[2], []byte(service)) {
		return fmt.Errorf("Expected service to be '%s', got '%s'", service, string(frames[2]))
	}

	return nil
}

func validateWorkerHeader(frames [][]byte) error {
	if length := len(frames[0]); length != 0 {
		return fmt.Errorf("First frame must be empty, got %v", string(frames[0]))
	}

	if !bytes.Equal(frames[1], mdp_WORKER) {
		return fmt.Errorf("Second frame must be %s, got %s", string(mdp_WORKER), string(frames[2]))
	}

	return nil
}

func validateWorkerReady(frames [][]byte, service string) error {
	//for i, x := range frames { fmt.Printf("%d: %v (%s)\n", i, x, string(x)) }
	if count := len(frames); count != 4 {
		return fmt.Errorf("Expected exactly 4 frames, got %d", count)
	}

	if err := validateWorkerHeader(frames); err != nil {
		return err
	}

	if len(frames[2]) != 1 || command(frames[2][0]) != c_READY {
		return fmt.Errorf("Third frame must be 0x%02x (READY cmd).", c_READY)
	}

	if !bytes.Equal(frames[3], []byte(service)) {
		return fmt.Errorf("Service mismatch. Expected %s but received %s", service, string(frames[4]))
	}
	return nil
}

func validateWorkerHeartbeat(frames [][]byte) error {
	// SPEC: Frame 0: Empty frame
	// SPEC: Frame 1: "MDPW01" (six bytes, representing MDP/Worker v0.1)
	// SPEC: Frame 2: 0x04 (one byte, representing HEARTBEAT)

	//for i, x := range frames { fmt.Printf("%d: %v\n", i, x) }

	if count := len(frames); count != 3 {
		return fmt.Errorf("Expected exactly 3 frames, got %d", count)
	}

	if err := validateWorkerHeader(frames); err != nil {
		return err
	}

	if len(frames[2]) != 1 || command(frames[2][0]) != c_HEARTBEAT {
		return fmt.Errorf("Third frame must be 0x%02x (HEARTBEAT cmd).", c_HEARTBEAT)
	}
	return nil
}

func validateWorkerRequest(frames [][]byte) error {
	if count := len(frames); count < 6 {
		return fmt.Errorf("Expected at least 6 frames, got %d", count)
	}

	if err := validateWorkerHeader(frames); err != nil {
		return err
	}

	if len(frames[2]) != 1 || command(frames[2][0]) != c_REQUEST {
		return fmt.Errorf("Third frame must be 0x%02x (REQUEST cmd).", c_REQUEST)
	}

	// I think the client address should probably be non-empty, though the spec doesn't say this. The spec says:
	// SPEC: The REQUEST and the REPLY cmds MUST contain precisely one client address frame.
	if len(frames[3]) == 0 {
		return fmt.Errorf("Fourth frame (client address) must be non-empty.")
	}

	if len(frames[4]) != 0 {
		return fmt.Errorf("Fifth frame must be empty")
	}
	return nil
}

func validateWorkerReply(frames [][]byte, client []byte) error {
	if count := len(frames); count < 6 {
		return fmt.Errorf("Expected at least 6 frames, got %d", count)
	}

	if err := validateWorkerHeader(frames); err != nil {
		return err
	}

	if len(frames[2]) != 1 || command(frames[2][0]) != c_REPLY {
		return fmt.Errorf("Third frame must be 0x%02x (REPLY cmd).", c_REPLY)
	}

	// I think the client address should probably be non-empty, though the spec doesn't say this. The spec says:
	// SPEC: The REQUEST and the REPLY cmds MUST contain precisely one client address frame.
	if len(frames[3]) == 0 {
		return fmt.Errorf("Fourth frame (client address) must be non-empty.")
	}

	if !bytes.Equal(frames[3], client) {
		return fmt.Errorf("Client address must match")
	}

	if len(frames[4]) != 0 {
		return fmt.Errorf("Fifth frame must be empty")
	}
	return nil
}

func durationInMilliseconds(d time.Duration) int {
	return int(d / time.Millisecond)
}

func czmqPollerSafeWait(poller *czmq.Poller, timeout_milliseconds int) *czmq.Sock {
	defer func() {
		if r := recover(); r != nil {
			if r == czmq.WaitAfterDestroyPanicMessage {
				// ignore
			} else {
				panic(r)
			}
		}
	}()
	return poller.Wait(timeout_milliseconds)
}
