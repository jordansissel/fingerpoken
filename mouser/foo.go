package main

import (
  "github.com/BurntSushi/xgb/xtest"
  "github.com/BurntSushi/xgb"
  "github.com/BurntSushi/xgb/xproto"
  "github.com/jordansissel/fingerpoken/util"
  "fmt"
  consul "github.com/hashicorp/consul/api"
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
  fmt.Printf("%#v\n", args)
  screen := xproto.Setup(m.x11).DefaultScreen(m.x11)
  cookie := xtest.FakeInputChecked(m.x11, xproto.MotionNotify, 0, 0, screen.Root, int16(args.X), int16(args.Y), 0)
  if cookie.Check() != nil {
    fmt.Println("FakeInput failed")
  }

  return nil
}

func main() {
  client, _ := consul.NewClient(consul.DefaultConfig())
  zj, err := util.NewZJServer() 
  if err != nil {
    fmt.Printf("NewZJRPC failure: %s\n", err)
    panic("!")
  }
  zj.RegisterWithConsul(client)
  zj.Register(NewMouse())
  err = zj.Loop()
  fmt.Printf("Loop: %s\n", err)
}
