package osdp

type OSDPTransceiver interface {
	Transmit(payload []byte) error
	Receive() ([]byte, error) // Byte slice received must only have bytes received returned, no extra padding
}
