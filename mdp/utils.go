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
		return fmt.Errorf("First frame must be empty")
	}

	if !bytes.Equal(frames[1], MDP_WORKER) {
		return fmt.Errorf("Second frame must be %s, got %s", string(MDP_WORKER), string(frames[2]))
	}

	return nil
}

func validateWorkerReady(frames [][]byte, service string) error {
	// Since this is a ROUTER socket, the first frame is the client id, so the MDP frames start at frames[1]
	//for i, x := range frames {
	//fmt.Printf("%d: %v\n", i, string(x))
	//}
	if count := len(frames); count < 4 {
		return fmt.Errorf("Expected at least 4 frames, got %d", count)
	}

	if err := validateWorkerHeader(frames); err != nil {
		return err
	}

	if !bytes.Equal(frames[2], []byte{C_READY}) {
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

	if !bytes.Equal(frames[2], []byte{C_HEARTBEAT}) {
		return fmt.Errorf("Third frame must be 0x%02x (HEARTBEAT command).", C_HEARTBEAT)
	}
	return nil
}
