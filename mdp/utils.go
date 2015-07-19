package mdp

import (
	"bytes"
	"fmt"
	"math/rand"
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

	if !bytes.Equal(frames[1], MDP_CLIENT) {
		return fmt.Errorf("Second frame must be %s, got %s", string(MDP_CLIENT), string(frames[2]))
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

	if !bytes.Equal(frames[1], MDP_WORKER) {
		return fmt.Errorf("Second frame must be %s, got %s", string(MDP_WORKER), string(frames[2]))
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

	if len(frames[2]) != 1 || Command(frames[2][0]) != C_READY {
		return fmt.Errorf("Third frame must be 0x%02x (READY command).", C_READY)
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

	if len(frames[2]) != 1 || Command(frames[2][0]) != C_HEARTBEAT {
		return fmt.Errorf("Third frame must be 0x%02x (HEARTBEAT command).", C_HEARTBEAT)
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

	if len(frames[2]) != 1 || Command(frames[2][0]) != C_REQUEST {
		return fmt.Errorf("Third frame must be 0x%02x (REQUEST command).", C_REQUEST)
	}

	// I think the client address should probably be non-empty, though the spec doesn't say this. The spec says:
	// SPEC: The REQUEST and the REPLY commands MUST contain precisely one client address frame.
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

	if len(frames[2]) != 1 || Command(frames[2][0]) != C_REPLY {
		return fmt.Errorf("Third frame must be 0x%02x (REPLY command).", C_REPLY)
	}

	// I think the client address should probably be non-empty, though the spec doesn't say this. The spec says:
	// SPEC: The REQUEST and the REPLY commands MUST contain precisely one client address frame.
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
