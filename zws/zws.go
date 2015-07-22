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
	"github.com/gorilla/websocket"
	//"github.com/gorilla/mux"
	"fmt"
	czmq "github.com/zeromq/goczmq"
	"net/http"
)

type ZWS struct {
	ws          *websocket.Conn
	zmq         *czmq.Sock
	socket_type SocketType
}

func NewZWS(endpoint_prefix string, w http.ResponseWriter, r *http.Request) (*ZWS, error) {
	query := r.URL.Query()
	socket_type_name, ok := query["type"]
	if !ok {
		return nil, &ZWSRequestMissingTypeParameterError{r.URL}
	}

	number, err := ParseSocketType(socket_type_name[0])
	if err != nil {
		return nil, &ZWSRequestInvalidTypeParameterError{socket_type_name[0], err}
	}
	socket_type := SocketType(number)

	if !socket_type.isValid() {
		return nil, &ZWSRequestInvalidTypeParameterError{socket_type_name[0], err}
	}

	endpoint_suffix, err := socket_type.EndpointSuffix()
	if err != nil {
		return nil, err
	}

	zmq, err := socket_type.Create(fmt.Sprintf("%s-%s", endpoint_prefix, endpoint_suffix))
	if err != nil {
		return nil, &ZeroMQSocketCreationError{err}
	}

	ws, err := upgrade(w, r)
	if err != nil {
		return nil, err
	}

	zws := &ZWS{
		ws:          ws,
		socket_type: socket_type,
		zmq:         zmq,
	}
	return zws, nil
}
