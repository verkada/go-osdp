package osdp

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

type OSDPMessage struct {
	MessageCode       OSDPCode
	PeripheralAddress byte
	MessageData       []byte
	SequenceNumber    byte
	secure            bool
	secureBlockType   byte
	secureBlockData   []byte
	Retries           uint32
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

func PacketFromMessage(osdpMessage *OSDPMessage) (*OSDPPacket, error) {

	if osdpMessage.secure {
		osdpPacket, err := NewSecurePacket(osdpMessage.MessageCode, osdpMessage.PeripheralAddress, osdpMessage.MessageData, osdpMessage.secureBlockType, osdpMessage.secureBlockData, osdpMessage.SequenceNumber, true)
		if err != nil {
			return nil, err
		}
		return osdpPacket, nil
	}
	osdpPacket, err := NewPacket(osdpMessage.MessageCode, osdpMessage.PeripheralAddress, osdpMessage.MessageData, osdpMessage.SequenceNumber, true)
	if err != nil {
		return nil, err
	}
	return osdpPacket, nil
}

func GenerateMAC(osdpMessage *OSDPMessage, IVC, SMAC1, SMAC2 []byte) ([]byte, error) {

	if !osdpMessage.secure {
		return nil, errors.New("Can only generate MAC for secure message")
	}

	osdpPacket, err := PacketFromMessage(osdpMessage)
	if err != nil {
		return nil, err
	}

	if len(IVC) != 16 {
		return nil, errors.New("Invalid IVC Length")
	}

	if len(SMAC1) != 16 {
		return nil, errors.New("Invalid SMAC 1 Length")
	}

	if len(SMAC2) != 16 {
		return nil, errors.New("Invalid SMAC 2 Length")
	}
	//TODO: Remove MAC, REMOVE CRC
	MAC := make([]byte, 16)
	osdpPacketBytes := osdpPacket.ToBytes()
	packetLength := len(osdpPacketBytes)

	if packetLength > 16 {
		block, err := aes.NewCipher(SMAC1)
		if err != nil {
			return nil, err
		}
		mode := cipher.NewCBCEncrypter(block, IVC)
		mode.CryptBlocks(MAC, osdpPacketBytes[0:packetLength-16])
		osdpPacketBytes = osdpPacketBytes[packetLength-16:]
	}

	// fmt.Println("Hello, playground")
	// osdpPacket := []byte{0x53, 0x3d, 0x0e, 0x00, 0x0e, 0x02, 0x15, 0x60, 0x4b, 0x56, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00}
	// IVC:= []byte{0x61, 0x6b, 0xb9, 0xf8, 0x1a, 0x0f, 0xfd, 0x62, 0x8d, 0x27, 0x45, 0x39, 0xb4, 0xd1, 0x0a, 0x86}
	// key := []byte{ 0xb3 ,0x63 ,0xdf ,0x85 ,0x12 ,0x13 ,0xde ,0x3a  ,0x4b, 0x50, 0xd9, 0x02 ,0x9a ,0x97, 0x3a, 0xd7}
	// block, err := aes.NewCipher(key)
	// if err != nil {

	// 	fmt.Println("Error New Cipher")
	// }
	// dst := make([]byte, 16)
	// mode := cipher.NewCBCEncrypter(block, IVC)
	// mode.CryptBlocks(dst, osdpPacket)
	// fmt.Println("MAC:", hex.Dump(dst))
}
