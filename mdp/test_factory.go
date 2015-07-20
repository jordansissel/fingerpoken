package mdp

type HelloGreeter struct{}

var helloGreeting = []byte("Nice to meet you!")

func (h *HelloGreeter) Request(request [][]byte) (response [][]byte, err error) {
	response = append(response, helloGreeting)
	return
}

func (h *HelloGreeter) Heartbeat()  {}
func (h *HelloGreeter) Disconnect() {}
