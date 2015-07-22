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
