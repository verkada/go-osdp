package osdp

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
)

type OSDPMessage struct {
	MessageCode       OSDPCode
	PeripheralAddress byte
	MessageData       []byte
	SequenceNumber    byte
	Secure            bool
	SecureBlockType   byte
	SecureBlockData   []byte
	Retries           uint32
	MAC               []byte
}

func NewOSDPMessage(osdpCode OSDPCode, peripheralAddress byte, sequenceNumber byte, msgData []byte) (*OSDPMessage, error) {
	if sequenceNumber < 0x00 || sequenceNumber > 0x03 {
		return nil, InvalidSequenceNumber
	}
	return &OSDPMessage{MessageCode: osdpCode, PeripheralAddress: peripheralAddress, MessageData: msgData, SequenceNumber: sequenceNumber, Secure: false}, nil
}

func NewSecureOSDPMessage(osdpCode OSDPCode, peripheralAddress byte, sequenceNumber byte, secureBlockType byte, secureBlockData []byte, msgData []byte) (*OSDPMessage, error) {
	if sequenceNumber < 0x00 || sequenceNumber > 0x03 {
		return nil, InvalidSequenceNumber
	}
	if secureBlockType < SCS_11 || secureBlockType > SCS_18 {
		return nil, InvalidSecureBlockType
	}
	return &OSDPMessage{MessageCode: osdpCode, PeripheralAddress: peripheralAddress, MessageData: msgData, SequenceNumber: sequenceNumber, Secure: true, SecureBlockType: secureBlockType, SecureBlockData: secureBlockData}, nil
}

func (osdpMessage *OSDPMessage) PacketFromMessage() (*OSDPPacket, error) {

	if osdpMessage.Secure {
		osdpPacket, err := NewSecurePacket(osdpMessage.MessageCode, osdpMessage.PeripheralAddress, osdpMessage.MessageData, osdpMessage.SecureBlockType, osdpMessage.SecureBlockData, osdpMessage.SequenceNumber, true)
		if err != nil {
			return nil, err
		}
		if osdpMessage.MAC != nil {
			osdpPacket.msgAuthenticationCode = make([]byte, 4)
			copy(osdpPacket.msgAuthenticationCode, osdpMessage.MAC[:4])
		}
		return osdpPacket, nil
	}

	osdpPacket, err := NewPacket(osdpMessage.MessageCode, osdpMessage.PeripheralAddress, osdpMessage.MessageData, osdpMessage.SequenceNumber, true)
	if err != nil {
		return nil, err
	}
	return osdpPacket, nil
}

func (osdpMessage *OSDPMessage) GenerateMAC(IVC, SMAC1, SMAC2 []byte) ([]byte, error) {

	if !osdpMessage.Secure {
		return nil, errors.New("Can only generate MAC for secure message")
	}

	osdpPacket, err := osdpMessage.PacketFromMessage()
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
	osdpPacketBytes = osdpPacketBytes[0 : len(osdpPacketBytes)-2] // Remove the CRC
	if osdpPacket.useMAC {
		osdpPacketBytes = osdpPacketBytes[0 : len(osdpPacketBytes)-4]
	}
	packetLength := len(osdpPacketBytes)
	// Apply Padding
	paddingRequired := 16 - (packetLength % 16)
	if packetLength < 16 {
		paddingRequired = 16 - packetLength
	}
	osdpPacketBytes = append(osdpPacketBytes, 0x80)
	restPadding := make([]byte, paddingRequired-1)
	osdpPacketBytes = append(osdpPacketBytes, restPadding...)
	packetLength = len(osdpPacketBytes)
	if packetLength > 16 {
		block, err := aes.NewCipher(SMAC1)
		if err != nil {
			return nil, err
		}
		for packetLength > 16 {
			mode := cipher.NewCBCEncrypter(block, IVC)
			mode.CryptBlocks(MAC, osdpPacketBytes[0:16])
			osdpPacketBytes = osdpPacketBytes[16:]
			packetLength = len(osdpPacketBytes)
			if copy(IVC, MAC) != 16 {
				return nil, errors.New("Unable to copy MAC into IV")
			}
		}

	}
	block, err := aes.NewCipher(SMAC2)
	if err != nil {
		return nil, err
	}
	fmt.Println(len(MAC), " ", len(IVC), " ", len(osdpPacketBytes))
	mode := cipher.NewCBCEncrypter(block, IVC)
	mode.CryptBlocks(MAC, osdpPacketBytes)
	osdpMessage.MAC = MAC
	return MAC, nil
}
