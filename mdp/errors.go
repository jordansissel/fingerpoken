package mdp

type Error interface {
  error
}

type ProtocolError struct {
}
