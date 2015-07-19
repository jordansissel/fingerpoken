package mdp
import (
  "math/rand"
  "fmt"
  "bytes"
)

func randomHex() (value string) {
	var length = rand.Int31n(10) + 5
	for x := int32(0); x < length; x += 1 {
		value = fmt.Sprintf("%s%02x", value, rand.Int31n(256))
	}
	return
}

func validateReadyMessage(frames [][]byte, service string) error {
  // Since this is a ROUTER socket, the first frame is the client id, so the MDP frames start at frames[1]
  //for i, x := range frames {
    //fmt.Printf("%d: %v\n", i, string(x))
  //}
  if count := len(frames); count < 5 {
    return fmt.Errorf("Expected at least 5 frames, got %d", count)
  }

  if length := len(frames[1]); length != 0 {
    return fmt.Errorf("First frame must be empty")
  }

  if !bytes.Equal(frames[2], MDP_WORKER) {
    return fmt.Errorf("Second frame must be %s, got %s", string(MDP_WORKER), string(frames[2]))
  }

  if !bytes.Equal(frames[3], []byte{C_READY}) {
    return fmt.Errorf("Third frame must be 0x%02x (READY command).", C_READY)
  }

  if !bytes.Equal(frames[4], []byte(service)) {
    return fmt.Errorf("Service mismatch. Expected %s but received %s", service, string(frames[4]))
  }
  return nil
}

func validateHeartbeat(frames [][]byte) error {
  // SPEC: Frame 0: Empty frame
  // SPEC: Frame 1: "MDPW01" (six bytes, representing MDP/Worker v0.1)
  // SPEC: Frame 2: 0x04 (one byte, representing HEARTBEAT)

  //for i, x := range frames { fmt.Printf("%d: %v\n", i, x) }
  
  // Since this is a ROUTER socket, the first frame is the client id, so the MDP frames start at frames[1]

  if count := len(frames); count != 4 {
    return fmt.Errorf("Expected exactly 4 frames, got %d", count)
  }

  if length := len(frames[1]); length != 0 {
    return fmt.Errorf("First frame must be empty")
  }

  if !bytes.Equal(frames[2], MDP_WORKER) {
    return fmt.Errorf("Second frame must be %s, got %s", string(MDP_WORKER), string(frames[2]))
  }

  if !bytes.Equal(frames[3], []byte{C_HEARTBEAT}) {
    return fmt.Errorf("Third frame must be 0x%02x (HEARTBEAT command).", C_HEARTBEAT)
  }
  return nil
}
