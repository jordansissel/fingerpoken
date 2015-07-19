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
}
