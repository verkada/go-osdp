package osdp

type OSDPMessage struct {
	MessageCode       OSDPCode
	PeripheralAddress byte
	MessageData       []byte
	SequenceNumber    byte
	secure            bool
	secureBlockType   byte
	secureBlockData   []byte
}

func NewOSDPMessage(osdpCode OSDPCode, peripheralAddress byte, sequenceNumber byte, msgData []byte) (*OSDPMessage, error) {
	if sequenceNumber < 0x00 || sequenceNumber > 0x03 {
		return nil, InvalidSequenceNumber
	}
	return &OSDPMessage{MessageCode: osdpCode, PeripheralAddress: peripheralAddress, MessageData: msgData, SequenceNumber: sequenceNumber, secure: false}, nil
}

func NewSecureOSDPMessage(osdpCode OSDPCode, peripheralAddress byte, sequenceNumber byte, secureBlockType byte, secureBlockData []byte, msgData []byte) (*OSDPMessage, error) {
	if sequenceNumber < 0x00 || sequenceNumber > 0x03 {
		return nil, InvalidSequenceNumber
	}
	if secureBlockType < SCS_11 || secureBlockType > SCS_18 {
		return nil, InvalidSecureBlockType
	}
	return &OSDPMessage{MessageCode: osdpCode, PeripheralAddress: peripheralAddress, MessageData: msgData, SequenceNumber: sequenceNumber, secure: true, secureBlockType: secureBlockType, secureBlockData: secureBlockData}, nil
}
