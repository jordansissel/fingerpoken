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
package zws

import (
	"fmt"
	"net/url"
)

type InvalidSocketTypeError struct {
	SocketType SocketType
}

func (e *InvalidSocketTypeError) Error() string {
	return fmt.Sprintf("Invalid socket type given: %s", e.SocketType)
}

type InvalidMessageTypeError struct {
	MessageType int
}

func (e *InvalidMessageTypeError) Error() string {
	return fmt.Sprintf("Unexpected message type: %d", e.MessageType)
}

type InvalidMoreFlagError struct {
	Flag byte
}

func (e *InvalidMoreFlagError) Error() string {
	return fmt.Sprintf("Invalid MORE flag: %c", e.Flag)
}

type ZeroMQSocketCreationError struct {
	Cause error
}

func (e *ZeroMQSocketCreationError) Error() string {
	return fmt.Sprintf("zeromq socket creation failed: %s", e.Cause)
}

type ZWSRequestMissingTypeParameterError struct {
	Url *url.URL
}

func (e *ZWSRequestMissingTypeParameterError) Error() string {
	return fmt.Sprintf("Request is missing `?type=SOCKET_TYPE` in path: %s", e.Url)
}

type ZWSRequestInvalidTypeParameterError struct {
	SocketTypeName string
	Cause          error
}

func (e *ZWSRequestInvalidTypeParameterError) Error() string {
	return fmt.Sprintf("`type=%s` error: %s", e.SocketTypeName, e.Cause)
}
