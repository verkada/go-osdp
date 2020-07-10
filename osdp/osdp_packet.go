package osdp

type OSDPPacket struct {
	startOfMessage            byte
	peripheralAddress         byte
	lsbLength                 byte
	msbLength                 byte
	msgCtrlInfo               byte
	securityBlockLength       byte
	securityBlockType         byte
	securityBlockData         []byte
	msgCode                   byte // OSDP Command or Reply
	messageData               []byte
	messageAuthenticationCode []byte //Max len 4
	lsbCheckSum               byte
	msbCheckSum               byte
}

const (
	OSDPSOM byte = 0x53
)

func NewOSDPPacket(msgCode byte, peripheralAddress byte, messageData []byte) (*OSDPPacket, error) {

	osdpPacket := &OSDPPacket{startOfMessage: OSDPSOM,
		peripheralAddress: peripheralAddress, msgCode: msgCode, messageData: messageData}
	//TODO: check that arguments meet OSDP spec
	return osdpPacket, nil
}

func (osdpPacket OSDPPacket) ToBytes() []byte {
	// TODO: implement
	return []byte{}
}
