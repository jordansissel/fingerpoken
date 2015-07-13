package util

import (
  "net/url"
  "net"
  "strconv"
  consul "github.com/hashicorp/consul/api"
)

func registerConsulService(client *consul.Client, name, endpoint string) error {
  agent := client.Agent()
  uri, _ := url.Parse(endpoint)
  _, port, _ := net.SplitHostPort(uri.Host)
  port_i, _ := strconv.ParseInt(port, 10, 32)

  service := &consul.AgentServiceRegistration{
    Name: "rpc",
    Port: int(port_i),
  }
  err := agent.ServiceRegister(service)
  return err
}

