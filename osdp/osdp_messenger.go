package osdp

import (
	"errors"
	"time"
)

type OSDPMessengerEvent int

// TODO: propagate event with error codes
const (
	OSDPDisconnected   OSDPMessengerEvent = 0
	OSDPConnected      OSDPMessengerEvent = 1
	OSDPReceiveTimeout OSDPMessengerEvent = 2
	OSDPReceiveError   OSDPMessengerEvent = 3
	OSDPTransmitError  OSDPMessengerEvent = 4
)

type OSDPMessenger struct {
	connected   bool
	transceiver OSDPTransceiver
}

type osdpResponseResult struct {
	responsePayload []byte
	responseErr     error
}

func NewOSDPMessenger(transceiver OSDPTransceiver) *OSDPMessenger {
	return &OSDPMessenger{connected: false, transceiver: transceiver}
}

func (osdpMessenger *OSDPMessenger) SendOSDPCommand(osdpMessage *OSDPMessage, timeout time.Duration) error {
	osdpPacket, err := NewPacket(osdpMessage.MessageCode, osdpMessage.PeripheralAddress, osdpMessage.MessageData, true)
	if err != nil {
		return err
	}
	return osdpMessenger.transceiver.Transmit(osdpPacket.ToBytes())
}

func (osdpMessenger *OSDPMessenger) ReceiveResponse(timeout time.Duration) (*OSDPMessage, error) {
	receiveChannel := make(chan osdpResponseResult, 1)
	go func() {
		responseData, err := osdpMessenger.transceiver.Receive() // TODO handle max length correctly
		responseResult := osdpResponseResult{responsePayload: responseData, responseErr: err}
		receiveChannel <- responseResult
	}()

	select {
	case response := <-receiveChannel:
		if response.responseErr != nil {
			return nil, response.responseErr
		}
		osdpPacket, err := NewPacketFromBytes(response.responsePayload)
		if err != nil {
			return nil, err
		}
		return &OSDPMessage{
			MessageCode:       OSDPCode(osdpPacket.msgCode),
			PeripheralAddress: osdpPacket.peripheralAddress, MessageData: osdpPacket.msgData,
		}, nil
	case <-time.After(timeout * time.Millisecond):
		return nil, errors.New("OSDPReceiveTimeout")
	}
}

func (osdpMessenger *OSDPMessenger) SendAndReceive(osdpMessage *OSDPMessage, writeTimeout time.Duration, readTimeout time.Duration) (*OSDPMessage, error) {
	err := osdpMessenger.SendOSDPCommand(osdpMessage, writeTimeout)
	if err != nil {
		return nil, err
	}
	osdpPacket, err := osdpMessenger.ReceiveResponse(readTimeout)
	if err != nil {
		return nil, err
	}
	return osdpPacket, nil
}
