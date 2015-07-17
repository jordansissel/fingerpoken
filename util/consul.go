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
