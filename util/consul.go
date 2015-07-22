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
package util

import (
	consul "github.com/hashicorp/consul/api"
	"net"
	"net/url"
	"strconv"
)

func ConsulRegisterService(client *consul.Client, name, endpoint string) (err error) {
	agent := client.Agent()
	uri, err := url.Parse(endpoint)
	if err != nil {
		return
	}
	_, port, err := net.SplitHostPort(uri.Host)
	if err != nil {
		return
	}
	port_i, err := strconv.ParseInt(port, 10, 32)
	if err != nil {
		return
	}

	service := &consul.AgentServiceRegistration{
		Name: name,
		Port: int(port_i),
	}
	err = agent.ServiceRegister(service)
	return
}
