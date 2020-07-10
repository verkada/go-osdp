package osdp

import (
	"errors"
	"time"
)

type OSDPConnectionEvent int

// TODO: propagate event with error codes
const (
	OSDPDisconnected   OSDPConnectionEvent = 0
	OSDPConnected      OSDPConnectionEvent = 1
	OSDPReceiveTimeout OSDPConnectionEvent = 2
	OSDPReceiveError   OSDPConnectionEvent = 3
	OSDPTransmitError  OSDPConnectionEvent = 4
)

type OSDPConnection struct {
	connected   bool
	transceiver OSDPTransceiver
}

type osdpResponseResult struct {
	responsePayload []byte
	responseErr     error
}

func NewOSDPConnection(transceiver OSDPTransceiver) *OSDPConnection {

	return &OSDPConnection{connected: false, transceiver: transceiver}
}

func (osdpConnection *OSDPConnection) SendOSDPCommand(osdpMessage *OSDPMessage) error {

	osdpPacket, err := NewOSDPPacket(byte(osdpMessage.osdpCommand), osdpMessage.peripheralAddress, osdpMessage.messageData)
	if err != nil {
		return err
	}
	return osdpConnection.transceiver.Transmit(osdpPacket.ToBytes())
}

func (osdpConnection *OSDPConnection) ReceiveResponse(timeout time.Duration) (*OSDPMessage, error) {

	receiveChannel := make(chan osdpResponseResult, 1)
	go func() {
		responseData, err := osdpConnection.transceiver.Receive(255) // TODO handle max length correctly
		responseResult := osdpResponseResult{responsePayload: responseData, responseErr: err}
		receiveChannel <- responseResult
	}()

	select {
	case response := <-receiveChannel:
		if response.responseErr != nil {
			return nil, response.responseErr
		}
	case <-time.After(timeout * time.Millisecond):
		return nil, errors.New("OSDPReceiveTimeout")
	}
	// TODO convert byte response into osdp Message
	return &OSDPMessage{}, nil
}

func (osdpConnection *OSDPConnection) SendAndReceive(osdpMessage *OSDPMessage, timeout time.Duration) (*OSDPMessage, error) {
	err := osdpConnection.SendOSDPCommand(osdpMessage)
	if err != nil {
		return nil, err
	}
	osdpPacket, err := osdpConnection.ReceiveResponse(timeout)
	if err != nil {
		return nil, err
	}
	return osdpPacket, nil
}
