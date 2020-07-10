package osdp

type OSDPTransceiver interface {
	Transmit(payload []byte) error
	Receive(max_length int) ([]byte, error)
}
