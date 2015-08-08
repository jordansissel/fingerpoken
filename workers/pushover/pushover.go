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
	log "github.com/Sirupsen/logrus"
	"github.com/thorduri/pushover"
)

type Pushover struct {
	client string
	user   string
}

type PushoverMessage struct {
	Message   string `json:message`
	Title     string `json:title`
	URL       string `json:url`
	URLTitle  string `json:url_title`
	Priority  int    `json:priority`
	Retry     int    `json:retry`
	Timestamp int64  `json:timestamp`

	//Device string
	//Sound string
}

type PushoverResponse struct {
	Status  int      `json:status`
	Request string   `json:request`
	Errors  []string `json:errors,omitempty`
}

type Empty struct{}

func (p *Pushover) Send(args *PushoverMessage, reply *PushoverResponse) error {
	papi, err := pushover.NewPushover(p.client, p.user)
	if err != nil {
		log.Printf("pushover.NewPushover error: %s", err)
		return err
	}

	message := pushover.Message{
		Message:  args.Message,
		Priority: args.Priority,
		Retry:    args.Retry,
	}

	if len(args.Title) > 0 {
		message.Title = args.Title
	}

	if len(args.URL) > 0 {
		message.Url = args.URL
	}

	if len(args.URLTitle) > 0 {
		message.UrlTitle = args.URLTitle
	}

	if args.Timestamp > 0 {
		message.Timestamp = args.Timestamp
	}

	// TODO(sissel): message.Sound
	// TODO(sissel): message.Device
	// TODO(sissel): message.Retry
	// TODO(sissel): message.Expire
	//
	request, _, err := papi.Push(&message)
	if err != nil {
		log.Printf("papi.Push() failed: %s", err)
		return err
	}

	reply.Request = request
	// TODO(sissel): What to do with receipt?
	return nil
}
