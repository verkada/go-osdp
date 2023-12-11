package osdp_application

type ExampleTransceiver struct{}

func (transceiver *ExampleTransceiver) Transmit(payload []byte) error {
	return nil
}

func (transceiver *ExampleTransceiver) Receive() ([]byte, error) {
	return []byte{0x53, 0x00, 0x08, 0x00, 0x04, 0x40, 0x89, 0x8E}, nil
}

func NewTransceiver() *ExampleTransceiver {
	return &ExampleTransceiver{}
}
