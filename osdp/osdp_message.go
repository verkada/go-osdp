package osdp

type OSDPMessage struct {
	MessageCode       OSDPCode
	PeripheralAddress byte
	MessageData       []byte
}

func NewOSDPMessage(osdpCode OSDPCode, peripheralAddress byte, msgData []byte) *OSDPMessage {

	return &OSDPMessage{MessageCode: osdpCode, PeripheralAddress: peripheralAddress, MessageData: msgData}
}
