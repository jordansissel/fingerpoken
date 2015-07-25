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
package zap

import (
	"fmt"
  "log"
	czmq "github.com/zeromq/goczmq"
  "testing"
  "math/rand"
)

func randomHex() (value string) {
	var length = rand.Int31n(10) + 5
	for x := int32(0); x < length; x += 1 {
		value = fmt.Sprintf("%s%02x", value, rand.Int31n(256))
	}
	return
}

type OpenAccess struct { }
func (o OpenAccess) Authorize(authRequest ZapRequest) (allow bool, err error) {
  return true, nil
}

//var agent *ZapAgent

//func init() {
  //agent, _ = NewZapAgent()
  //go agent.Run(&OpenAccess{})
  ////defer agent.Destroy()
//}


func TestZAPAllow(t *testing.T) {
  agent, _ := NewZapAgent()
  go agent.Run(&OpenAccess{})
  defer agent.Destroy()

  domain := randomHex()
  text := randomHex()

  server := czmq.NewSock(czmq.Pull)
  defer server.Destroy()
  server.SetZapDomain(domain)
  port, err := server.Bind("tcp://127.0.0.1:*")
  if err != nil {
    t.Errorf("Binding server failed; %s", err)
    return
  }

  client := czmq.NewSock(czmq.Push)
  defer client.Destroy()
  client.Connect(fmt.Sprintf("tcp://127.0.0.1:%d", port))
  err = client.SendFrame([]byte(text), 0)
  if err != nil {
    t.Errorf("client.SendFrame failed: %s", err)
    return
  }

  output, _, err := server.RecvFrame()
  if err != nil {
    t.Errorf("server.RecvFrame() failed: %s", err)
  }

  if string(output) != text {
    t.Errorf("Output did not match.")
    return
  }
}

type DenyAccess struct { }
func (d DenyAccess) Authorize(authRequest ZapRequest) (allow bool, err error) {
  log.Printf("DENIED")
  return false, nil
}

func TestZAPDeny(t *testing.T) {
  agent, _ := NewZapAgent()
  go agent.Run(&DenyAccess{})
  defer agent.Destroy()

  domain := randomHex()
  text := randomHex()

  server := czmq.NewSock(czmq.Pull)
  defer server.Destroy()
  server.SetZapDomain(domain)
  port, err := server.Bind("tcp://127.0.0.1:*")
  if err != nil {
    t.Errorf("Binding server failed; %s", err)
    return
  }

  client := czmq.NewSock(czmq.Push)
  defer client.Destroy()
  err = client.Connect(fmt.Sprintf("tcp://127.0.0.1:%d", port))
  if err != nil {
    t.Errorf("client.Connect failed: %s", err)
    return
  }

  err = client.SendFrame([]byte(text), 0)
  if err != nil {
    t.Errorf("client.SendFrame failed: %s", err)
    return
  }

  output, _, err := server.RecvFrame()
  if err != nil {
    t.Errorf("server.RecvFrame() failed: %s", err)
  }

  if string(output) != text {
    t.Errorf("Output did not match.")
    return
  }


}

