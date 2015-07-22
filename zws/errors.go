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
