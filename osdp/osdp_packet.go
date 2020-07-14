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
	OSDPSOM                     byte   = 0x53
	minPeripheralAddress        byte   = 0x00
	maxPeripheralAddress        byte   = 0x7F
	msgControlChecksumMask      byte   = 0x04
	msgControlSecureMask        byte   = 0x08
	minimumPacketLengthUnsecure uint16 = 8
)

func NewPacket(msgCode OSDPCode, peripheralAddress byte, msgData []byte, integrityCheck bool) (*OSDPPacket, error) {
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
	var messageLengthUint uint16 = minimumPacketLengthUnsecure + uint16(len(securityBlockData)+len(msgAuthenticationCode)+len(msgData))
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

func NewPacketFromBytes(payload []byte) (*OSDPPacket, error) {

	// Check that payload meets minimum OSDP spec size
	var payloadLength uint16 = uint16(len(payload))
	if payloadLength < minimumPacketLengthUnsecure {
		return nil, errors.New("Payload size less than minimum possible OSDP payload size")
	}

	currentIndex := 0
	// Check that start of message follows OSDP spec
	startOfMessage := OSDPSOM
	if payload[currentIndex] != startOfMessage {
		return nil, errors.New("Invalid OSDP SOM")
	}

	currentIndex++
	// Check that the peripheral Address is in range
	peripheralAddress := payload[currentIndex]
	if peripheralAddress < minPeripheralAddress || peripheralAddress > maxPeripheralAddress {
		return nil, errors.New("Peripheral Address out of range")
	}

	// Parse the message length
	currentIndex++
	var messageLength uint16 = uint16(payload[currentIndex] | (payload[currentIndex+1] << 4))
	bytesRemaining := messageLength - minimumPacketLengthUnsecure // TODO: Add more if secure

	// Check the message control info. TODO: Check for secure, MAC etc
	currentIndex += 2
	msgControlInfo := payload[currentIndex]
	integrityCheck := false
	if (msgControlInfo & msgControlChecksumMask) == msgControlChecksumMask {
		integrityCheck = true
	}

	// TODO check the security block length if secure

	currentIndex++
	// Check the message code
	msgCode := payload[currentIndex]
	currentIndex++
	//TODO: if MAC then subtract 4 from bytes remaining to get length of msgData
	msgData := payload[currentIndex : currentIndex+int(bytesRemaining)]

	currentIndex += int(bytesRemaining)

	lsbChecksum := payload[currentIndex]
	currentIndex++
	msbChecksum := payload[currentIndex]

	osdpPacket, err := NewPacket(OSDPCode(msgCode), peripheralAddress, msgData, integrityCheck)
	if err != nil {
		return nil, err
	}

	if lsbChecksum != osdpPacket.lsbChecksum || msbChecksum != osdpPacket.msbChecksum {
		return nil, errors.New("Packet failed to pass checksum")
	}

	return osdpPacket, err
}
