package osdp

import (
	"context"
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
	readDoneContext, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	payload := []byte{}
	for {
		select {
		default:
			responseData, err := osdpMessenger.transceiver.Receive() // TODO handle max length correctly
			if err != nil {
				return nil, err
			}
			payload = append(payload, responseData...)
			osdpPacket, err := NewPacketFromBytes(payload)
			if err == nil {
				return &OSDPMessage{
					MessageCode:       OSDPCode(osdpPacket.msgCode),
					PeripheralAddress: osdpPacket.peripheralAddress, MessageData: osdpPacket.msgData,
				}, nil
			}
			// Keep Receiving until we get a valid packet, timeout or error
			if err != PacketIncompleteError {
				return nil, err
			}

		case <-readDoneContext.Done():
			osdpPacket, err := NewPacketFromBytes(payload)
			if err == nil {
				return &OSDPMessage{
					MessageCode:       OSDPCode(osdpPacket.msgCode),
					PeripheralAddress: osdpPacket.peripheralAddress, MessageData: osdpPacket.msgData,
				}, nil
			}
			return nil, OSDPReceiveTimeoutError
		}
	}
}

func (osdpMessenger *OSDPMessenger) SendAndReceive(osdpMessage *OSDPMessage, writeTimeout time.Duration, readTimeout time.Duration) (*OSDPMessage, error) {
	err := osdpMessenger.SendOSDPCommand(osdpMessage, writeTimeout)
	if err != nil {
		return nil, err
	}
	return osdpMessenger.ReceiveResponse(readTimeout)
}
