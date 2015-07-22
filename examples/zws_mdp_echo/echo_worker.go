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
  "log"
)

type echoWorker struct{}

func (e *echoWorker) Request(request [][]byte) (response [][]byte, err error) {
	for i, x := range request { log.Printf("Echo worker: frame %d: %v (%s)\n", i, x, string(x)) }
	response = append(response, request...)
	return
}

func (e *echoWorker) Heartbeat()  {}
func (e *echoWorker) Disconnect() {}
