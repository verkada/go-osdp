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
	if osdpMessage.secure {
		osdpPacket, err := NewSecurePacket(osdpMessage.MessageCode, osdpMessage.PeripheralAddress, osdpMessage.MessageData, osdpMessage.secureBlockType, osdpMessage.secureBlockData, osdpMessage.SequenceNumber, true)
		if err != nil {
			return err
		}
		return osdpMessenger.transceiver.Transmit(osdpPacket.ToBytes())

	}
	osdpPacket, err := NewPacket(osdpMessage.MessageCode, osdpMessage.PeripheralAddress, osdpMessage.MessageData, osdpMessage.SequenceNumber, true)
	if err != nil {
		return err
	}
	return osdpMessenger.transceiver.Transmit(osdpPacket.ToBytes())
}

func (osdpMessenger *OSDPMessenger) ReceiveResponse(timeout time.Duration) (*OSDPMessage, error) {
	type osdpReceiveMessage struct {
		message *OSDPMessage
		err     error
	}
	osdpReceiveChan := make(chan osdpReceiveMessage, 1)
	go func() {
		payload := []byte{}
		for {
			responseData, err := osdpMessenger.transceiver.Receive()
			if err != nil {
				osdpReceiveChan <- osdpReceiveMessage{message: nil, err: err}
			}

			payload = append(payload, responseData...)
			osdpPacket, err := NewPacketFromBytes(payload)
			if err == nil {
				sequenceNumber := osdpPacket.msgCtrlInfo & 0x03
				osdpReceiveChan <- osdpReceiveMessage{message: &OSDPMessage{
					MessageCode:       OSDPCode(osdpPacket.msgCode),
					PeripheralAddress: osdpPacket.peripheralAddress, MessageData: osdpPacket.msgData,
					SequenceNumber: sequenceNumber,
				}, err: nil}
				break
			}
			// Keep Receiving until we get a valid packet, timeout or error
			if err != PacketIncompleteError && err != InvalidSOMError {
				osdpReceiveChan <- osdpReceiveMessage{message: nil, err: err}
			}
		}
	}()

	select {
	case msg := <-osdpReceiveChan:
		return msg.message, msg.err

	case <-time.After(timeout):
		err := osdpMessenger.transceiver.Reset()
		if err != nil {
			return nil, err
		}
		return nil, OSDPReceiveTimeoutError
	}
}

func (osdpMessenger *OSDPMessenger) SendAndReceive(osdpMessage *OSDPMessage, writeTimeout time.Duration, readTimeout time.Duration) (*OSDPMessage, error) {
	err := osdpMessenger.SendOSDPCommand(osdpMessage, writeTimeout)
	if err != nil {
		return nil, err
	}
	return osdpMessenger.ReceiveResponse(readTimeout)
}
