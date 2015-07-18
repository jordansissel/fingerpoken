package main

import (
  "testing"
  "time"
)

func Test_NotificationBroadcaster(t *testing.T) {
  nbc := NotificationBroadcaster{}

  if v := len(nbc.subscribers); v != 0 {
    t.Errorf("Default subscriber count must be 0, is %d", v)
  }

  c := make(chan *Notification)
  nbc.Subscribe(c)

  c2 := make(chan *Notification)
  nbc.Subscribe(c2)
  go func() { for { <- c2 } }()

  if v := len(nbc.subscribers); v != 2 {
    t.Errorf("subscriber count after two subscriptions must be 2, is %d", v)
  }

  go func() { 
    n := Notification{}
    nbc.Publish(&n) 
  }()

  select {
  case <- c:
    // OK
  case <- time.After(50 * time.Millisecond):
    t.Errorf("timeout waiting for a notification")
  }

  nbc.Unsubscribe(c)
  if v := len(nbc.subscribers); v != 1 {
    t.Errorf("subscriber count after two subscription and one unsubscription must be 1, is %d", v)
  }
}
