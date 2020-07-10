package osdp

type OSDPMessage struct {
	osdpCommand       OSDPCommand
	peripheralAddress byte
	messageData       []byte
}

func NewOSDPMessage(osdpCommand OSDPCommand, peripheralAddress byte, messageData []byte) *OSDPMessage {

	return &OSDPMessage{osdpCommand: osdpCommand, peripheralAddress: peripheralAddress, messageData: messageData}
}
