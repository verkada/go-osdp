package osdp

import (
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

func NewOSDPMessenger(transceiver OSDPTransceiver, secure bool) *OSDPMessenger {
	return &OSDPMessenger{connected: false, transceiver: transceiver}
}

func (osdpMessenger *OSDPMessenger) SendOSDPCommand(osdpMessage *OSDPMessage, timeout time.Duration) error {
	// TODO Implement write timeout
	osdpPacket, err := osdpMessage.PacketFromMessage()
	if err != nil {
		return err
	}

	return osdpMessenger.transceiver.Transmit(osdpPacket.ToBytes())
}

func (osdpMessenger *OSDPMessenger) ReceiveResponse(timeout time.Duration) (*OSDPMessage, error) {
	payload := []byte{}
	timeStart := time.Now()
	for {
		responseData, err := osdpMessenger.transceiver.Receive()
		if err != nil {
			return nil, err
		}

		payload = append(payload, responseData...)
		osdpPacket, err := NewPacketFromBytes(payload)
		if err == nil {
			sequenceNumber := osdpPacket.msgCtrlInfo & 0x03
			return &OSDPMessage{
				MessageCode:       OSDPCode(osdpPacket.msgCode),
				PeripheralAddress: osdpPacket.peripheralAddress, MessageData: osdpPacket.msgData,
				SequenceNumber:  sequenceNumber,
				MAC:             osdpPacket.msgAuthenticationCode,
				SecureBlockData: osdpPacket.securityBlockData,
				SecureBlockType: osdpPacket.securityBlockType,
				Secure:          osdpPacket.secure,
			}, nil

		}
		// Keep Receiving until we get a valid packet, timeout or error
		if err != PacketIncompleteError && err != InvalidSOMError {
			return nil, err
		}
		if time.Since(timeStart) > timeout {
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
