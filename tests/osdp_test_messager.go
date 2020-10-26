package main

type MockTransceiver struct {
	timesCalled int
}

func (transceiver *MockTransceiver) Transmit(payload []byte) error {
	return nil
}

func (transceiver *MockTransceiver) Receive() ([]byte, error) {
	timesCalled := transceiver.timesCalled
	payload := []byte{0x53, 0x00, 0x08, 0x00, 0x04, 0x40, 0x89, 0x8E}
	returnPayload := payload[timesCalled : timesCalled+1]
	transceiver.timesCalled = (transceiver.timesCalled + 1) % len(payload)
	return returnPayload, nil
}
