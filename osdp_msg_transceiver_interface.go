package osdp

type OSDPTransceiver interface {
	Transmit(payload []byte) error
	Receive() ([]byte, error)
}
