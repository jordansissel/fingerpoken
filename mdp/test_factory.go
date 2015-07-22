package mdp

type helloGreeter struct{}

var helloGreeting = []byte("Nice to meet you!")

func (h *helloGreeter) Request(request [][]byte) (response [][]byte, err error) {
	response = append(response, helloGreeting)
	return
}

func (h *helloGreeter) Heartbeat()  {}
func (h *helloGreeter) Disconnect() {}
