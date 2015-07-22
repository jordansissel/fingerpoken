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
  "sync"
  "log"
)

type NotificationBroadcaster struct {
  subscribers []chan *Notification
  mutex sync.Mutex
}

func (nbc *NotificationBroadcaster) Publish(n *Notification) {
  nbc.mutex.Lock()
  defer func() { nbc.mutex.Unlock() }()
  for _, c := range nbc.subscribers {
    notify(nbc, c, n)
  }
}

func notify(nbc *NotificationBroadcaster, c chan *Notification, n *Notification) {
  defer func() {
    if r := recover(); r != nil {
      // Panic because we wrote to a closed channel.
      // Unsubscribe this channel now.
      nbc.unsubscribeWithLock(c)
    }
  }()
  c <- n
}

func (nbc *NotificationBroadcaster) Subscribe(c chan *Notification) {
  nbc.mutex.Lock()
  defer func() { nbc.mutex.Unlock() }()
  nbc.subscribers = append(nbc.subscribers, c)
}

func (nbc *NotificationBroadcaster) Unsubscribe(c chan *Notification) {
  nbc.mutex.Lock()
  defer func() { nbc.mutex.Unlock() }()
  nbc.unsubscribeWithLock(c)
}


func (nbc *NotificationBroadcaster) unsubscribeWithLock(c chan *Notification) {
  var index int
  for i, n := range nbc.subscribers {
    if n == c {
      log.Printf("Unsub %v\n", c);
      index = i
      break
    }
  }

  nbc.subscribers = append(nbc.subscribers[:index], nbc.subscribers[index+1:]...)
}

