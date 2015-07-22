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

