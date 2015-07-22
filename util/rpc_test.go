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
package util

import (
	"encoding/json"
	"fmt"
	czmq "github.com/zeromq/goczmq"
	"math/rand"
	"strings"
	"testing"
)

func Test_ZJServer(t *testing.T) {
	zj, err := NewZJServer()
	if err != nil {
		t.Errorf("NewZJServer should not fail. Got (%v, #v)", zj, err)
	}

	if zj == nil {
		t.Errorf("NewZJServer should not return nil for first value. Got (%v, #v)", zj, err)
	}

	endpoint := "tcp://*:*"
	err = zj.Bind(endpoint)
	if err != nil {
		t.Errorf("ZJServer#Bind(%v) should not fail. Error: %v ", err)
	}

	if zj.endpoint == "" {
		t.Errorf("ZJServer#endpoint must not be empty after a successful #Bind call")
	}
}

type Foo struct {
	bar int
}

func (f *Foo) Bar(args interface{}, reply *int) (err error) {
	*reply = f.bar
	err = nil
	return
}

func Test_ZJServer_RPC_With_Good_Request(t *testing.T) {
	// Use a random number here to add a bit of entropy into our tests.
	foo := Foo{bar: rand.Int()}

	zj, _ := NewZJServer()
	zj.Register(&foo)
	zj.Bind("tcp://*:*")
	endpoint := zj.endpoint

	go zj.Once()

	req, _ := czmq.NewReq(endpoint)
	req.SendFrame([]byte(`{ "method": "Foo.Bar", "params": [ ], "id": 1 }`), 0)
	response, err := req.RecvMessage()
	if err != nil {
		t.Errorf("czmq.Sock.RecvMessage failed: %v\n", err)
	}

	actual := strings.Trim(string(response[0]), "\r\n")
	expected := fmt.Sprintf(`{"id":1,"result":%d,"error":null}`, foo.bar)
	if actual != expected {
		t.Errorf("Response was not expected:\n  Wanted: %#v\n     Got: %#v\n", expected, actual)
	}

}

func Test_ZJServer_RPC_With_Invalid_Method(t *testing.T) {
	foo := Foo{bar: rand.Int()}

	zj, _ := NewZJServer()
	zj.Register(&foo)
	zj.Bind("tcp://*:*")
	endpoint := zj.endpoint

	go zj.Once()

	req, _ := czmq.NewReq(endpoint)
	req.SendFrame([]byte(`{ "method": "Foo.NoSuchMethod", "params": [ ], "id": 1 }`), 0)
	response, err := req.RecvMessage()
	if err != nil {
		t.Errorf("czmq.Sock.RecvMessage failed: %v\n", err)
	}

	var obj map[string]interface{}
	json.Unmarshal(response[0], &obj)

	switch e := obj["error"].(type) {
	default:
		t.Errorf("Expected a string value for the error, got %v", obj["error"])
	case string:
		if e != "rpc: can't find method Foo.NoSuchMethod" {
			t.Errorf("Foo.NoSuchMethod call should have failed with expected error, but got: %v", e)
		}
	}
}
