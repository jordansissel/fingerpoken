package mdp

type HelloGreeter struct{}

var HELLO_GREETING = []byte("Nice to meet you!")

func (h *HelloGreeter) Request(request [][]byte) (response [][]byte, err error) {
	response = append(response, HELLO_GREETING)
	return
}

func (h *HelloGreeter) Heartbeat()  {}
func (h *HelloGreeter) Disconnect() {}
