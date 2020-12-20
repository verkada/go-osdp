package main

import "time"

type MockTransceiver struct {
	timesCalled int
}

type SlowTransceiver struct {
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

func (transceiver *SlowTransceiver) Transmit(payload []byte) error {
	return nil
}

func (transceiver *SlowTransceiver) Receive() ([]byte, error) {
	time.Sleep(300 * time.Millisecond)
	payload := []byte{0x53, 0x00, 0x08, 0x00, 0x04, 0x40, 0x89, 0x8E}
	return payload, nil
}

func (transceiver *SlowTransceiver) Reset() error {
	return nil
}

func (transceiver *MockTransceiver) Reset() error {
	return nil
}
