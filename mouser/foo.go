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
package main

import (
	"fmt"
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgb/xtest"
	"os"

	//consul "github.com/hashicorp/consul/api"
	//"github.com/jordansissel/fingerpoken/util"
	"github.com/jordansissel/fingerpoken/mdp"
)

type Mouse struct {
	x11 *xgb.Conn
}

func NewMouse() *Mouse {
	mouse := &Mouse{}

	x11, err := xgb.NewConn()
	if err != nil {
		panic("err")
	}

	xtest.Init(x11)

	mouse.x11 = x11
	return mouse
}

type MoveArgs struct {
	X int `json:x`
	Y int `json:y`
}

func (m *Mouse) Move(args *MoveArgs, reply *int) error {
	//fmt.Printf("%#v\n", args)
	screen := xproto.Setup(m.x11).DefaultScreen(m.x11)
	cookie := xtest.FakeInputChecked(m.x11, xproto.MotionNotify, 0, 0, screen.Root, int16(args.X), int16(args.Y), 0)
	if cookie.Check() != nil {
		fmt.Println("FakeInput failed")
	}

	return nil
}

func main() {
	broker_endpoint := os.Args[1]
	service := os.Args[2]
	w := mdp.NewJSONRPCWorker(broker_endpoint, service)
	w.Register(NewMouse())
	w.Run()

	//err = zj.Loop()
	//fmt.Printf("Loop: %s\n", err)
}
