package osdp

type OSDPMessage struct {
	MessageCode       OSDPCode
	PeripheralAddress byte
	MessageData       []byte
	SequenceNumber    byte
}

func NewOSDPMessage(osdpCode OSDPCode, peripheralAddress byte, sequenceNumber byte, msgData []byte) *OSDPMessage {
	return &OSDPMessage{MessageCode: osdpCode, PeripheralAddress: peripheralAddress, MessageData: msgData, SequenceNumber: sequenceNumber}
}
