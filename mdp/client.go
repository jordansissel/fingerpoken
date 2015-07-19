package mdp

import (
	"bytes"
	"fmt"
	czmq "github.com/zeromq/goczmq"
)

type Client struct {
	sock   *czmq.Sock
	broker string
}

func NewClient(broker string) (c *Client) {
	c = &Client{broker: broker}
	return
}

// TODO(sissel): With timeout?
func (c *Client) Send(service string, body []byte) (response *[]byte, err error) {
	err = c.ensure_connected()
	if err != nil {
		return
	}

	// Since we're using a REQ socket, we use a 3-frame message instead of the 4-frame message a DEALER would use.
	// TODO(sissel): The body can occupy more than 1 frame, let's maybe support that some day?
	var request [3][]byte = [3][]byte{
		MDP_CLIENT,
		[]byte(service),
		body,
	}
	err = c.sock.SendMessage(request[:])
	reply, err := c.sock.RecvMessage()
	if err != nil {
		return
	}

	if frames := len(reply); frames < 3 {
		err = fmt.Errorf("Majordomo protocol problem. A REPLY must be at least 3 frames, got %d frames in a message.", frames)
		return
	}

	if !bytes.Equal(reply[0], MDP_CLIENT) {
		err = fmt.Errorf("Majordomo protocol problem. Expected first frame to be `%s`. Got something else.", string(MDP_CLIENT))
		return
	}

	// Should we bother checking the `service` frame (reply[1]) ?
	response = &reply[2]
	return
}

func (c *Client) ensure_connected() error {
	if c.sock != nil {
		return nil
	}

	var err error
	c.sock, err = czmq.NewReq(c.broker)
	if err != nil {
		return err
	}
	return nil
}
