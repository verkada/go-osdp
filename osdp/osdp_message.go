package osdp

type OSDPMessage struct {
	osdpCode          OSDPCode
	peripheralAddress byte
	messageData       []byte
}

func NewOSDPMessage(osdpCode OSDPCode, peripheralAddress byte, messageData []byte) *OSDPMessage {

	return &OSDPMessage{osdpCode: osdpCode, peripheralAddress: peripheralAddress, messageData: messageData}
}
