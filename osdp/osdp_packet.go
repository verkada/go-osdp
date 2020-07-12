package osdp

import (
	"encoding/binary"
	"errors"

	"github.com/sigurn/crc16"
)

type OSDPPacket struct {
	startOfMessage        byte
	peripheralAddress     byte
	lsbLength             byte
	msbLength             byte
	msgCtrlInfo           byte   // TODO Support
	securityBlockLength   byte   // TODO Support
	securityBlockType     byte   // TODO Support
	securityBlockData     []byte // TODO Support
	msgCode               byte   // OSDP Command or Reply
	msgData               []byte
	msgAuthenticationCode []byte //Max len 4 // TODO support
	lsbChecksum           byte
	msbChecksum           byte
	secure                bool
	useMAC                bool
}

const (
	OSDPSOM                byte = 0x53
	minPeripheralAddress   byte = 0x00
	maxPeripheralAddress   byte = 0x7F
	msgControlChecksumMask byte = 0x04
	msgControlSecureMask   byte = 0x08
)

func NewOSDPPacket(msgCode OSDPCode, peripheralAddress byte, msgData []byte, integrityCheck bool) (*OSDPPacket, error) {
	//TODO: check that arguments meet OSDP spec, assert msgData is the right size
	if peripheralAddress < minPeripheralAddress || peripheralAddress > maxPeripheralAddress {
		return nil, errors.New("Peripheral Address out of range")
	}

	// TODO: Support sequence number
	var msgControlInfo byte = 0
	if integrityCheck == true {
		msgControlInfo |= msgControlChecksumMask
	}

	var msgAuthenticationCode []byte = []byte{}
	var securityBlockData []byte = []byte{}
	var minimumOSDPPacketLengthUnsecure uint16 = 8
	var messageLengthUint uint16 = minimumOSDPPacketLengthUnsecure + uint16(len(securityBlockData)+len(msgAuthenticationCode)+len(msgData))
	messageLength := make([]byte, 2)
	binary.LittleEndian.PutUint16(messageLength, messageLengthUint)
	osdpPacket := &OSDPPacket{startOfMessage: OSDPSOM,
		peripheralAddress: peripheralAddress, lsbLength: messageLength[0], msbLength: messageLength[1],
		msgCtrlInfo: msgControlInfo, securityBlockLength: 0x00, securityBlockType: 0x00, securityBlockData: nil,
		msgCode: byte(msgCode), msgData: msgData, msgAuthenticationCode: nil,
		lsbChecksum: 0x00, msbChecksum: 0x00, secure: false, useMAC: false}

	osdpPacketBytes := osdpPacket.ToBytes()
	packetBytesSizeWithoutChecksum := len(osdpPacketBytes) - 2
	crc16Table := crc16.MakeTable(crc16.CRC16_AUG_CCITT)
	checksumUint := crc16.Checksum(osdpPacketBytes[:packetBytesSizeWithoutChecksum], crc16Table)
	checksum := make([]byte, 2)
	binary.LittleEndian.PutUint16(checksum, checksumUint)
	osdpPacket.lsbChecksum = checksum[0]
	osdpPacket.msbChecksum = checksum[1]
	return osdpPacket, nil
}

func (osdpPacket *OSDPPacket) ToBytes() []byte {

	var packetBytes []byte = []byte{osdpPacket.startOfMessage, osdpPacket.peripheralAddress,
		osdpPacket.lsbLength, osdpPacket.msbLength, osdpPacket.msgCtrlInfo,
	}
	if osdpPacket.secure {
		packetBytes = append(packetBytes, osdpPacket.securityBlockLength, osdpPacket.securityBlockType)
		packetBytes = append(packetBytes, osdpPacket.securityBlockData...)
	}

	packetBytes = append(packetBytes, osdpPacket.msgCode)
	packetBytes = append(packetBytes, osdpPacket.msgData...)

	if osdpPacket.useMAC {
		packetBytes = append(packetBytes, osdpPacket.msgAuthenticationCode...)
	}

	packetBytes = append(packetBytes, osdpPacket.lsbChecksum, osdpPacket.msbChecksum)

	return packetBytes
}
