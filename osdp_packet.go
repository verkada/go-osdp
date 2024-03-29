package osdp

import (
	"encoding/binary"

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
	msgAuthenticationCode []byte // Max len 4 // TODO support
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
	maxSecureBlockLength        int    = 0xFE
)

func NewSecurePacket(msgCode OSDPCode, peripheralAddress byte, msgData []byte, secureBlockType byte, secureBlockData []byte, sequenceNumber byte, integrityCheck bool) (*OSDPPacket, error) {
	if (peripheralAddress&maxPeripheralAddress) < minPeripheralAddress || (peripheralAddress&maxPeripheralAddress) > maxPeripheralAddress {
		return nil, AddressOutOfRangeError
	}

	if sequenceNumber > 0x03 {
		return nil, InvalidSequenceNumber
	}

	var msgControlInfo byte = 0
	if integrityCheck == true {
		msgControlInfo |= msgControlChecksumMask
	}

	msgControlInfo |= msgControlSecureMask

	msgControlInfo |= sequenceNumber

	if len(secureBlockData) > maxSecureBlockLength {
		return nil, SecureBlockDataLengthError
	}
	secureBlockLength := make([]byte, 1)
	secureBlockPayloadLen := int8(len(secureBlockData))
	secureBlockLength[0] = 0x02 + byte(secureBlockPayloadLen)

	useMAC := false
	var msgAuthenticationCode []byte
	if secureBlockType > SCS_14 {
		msgAuthenticationCode = []byte{0, 0, 0, 0}
		useMAC = true
	}
	var messageLengthUint uint16 = minimumPacketLengthUnsecure + uint16(int(secureBlockLength[0])+len(msgAuthenticationCode)+len(msgData))
	messageLength := make([]byte, 2)
	binary.LittleEndian.PutUint16(messageLength, messageLengthUint)
	osdpPacket := &OSDPPacket{
		startOfMessage:    OSDPSOM,
		peripheralAddress: peripheralAddress, lsbLength: messageLength[0], msbLength: messageLength[1],
		msgCtrlInfo: msgControlInfo, securityBlockLength: secureBlockLength[0], securityBlockType: secureBlockType, securityBlockData: secureBlockData,
		msgCode: byte(msgCode), msgData: msgData, msgAuthenticationCode: msgAuthenticationCode,
		lsbChecksum: 0x00, msbChecksum: 0x00, secure: true, useMAC: useMAC,
	}
	osdpPacket.calculateCRC()
	return osdpPacket, nil
}

func NewPacket(msgCode OSDPCode, peripheralAddress byte, msgData []byte, sequenceNumber byte, integrityCheck bool) (*OSDPPacket, error) {
	// TODO: check that arguments meet OSDP spec, assert msgData is the right size
	if (peripheralAddress&maxPeripheralAddress) < minPeripheralAddress || (peripheralAddress&maxPeripheralAddress) > maxPeripheralAddress {
		return nil, AddressOutOfRangeError
	}

	if sequenceNumber > 0x03 {
		return nil, InvalidSequenceNumber
	}

	// TODO: Support sequence number
	var msgControlInfo byte = 0
	if integrityCheck == true {
		msgControlInfo |= msgControlChecksumMask
	}

	msgControlInfo |= sequenceNumber

	var msgAuthenticationCode []byte = []byte{}
	var securityBlockData []byte = []byte{}
	var messageLengthUint uint16 = minimumPacketLengthUnsecure + uint16(len(securityBlockData)+len(msgAuthenticationCode)+len(msgData))
	messageLength := make([]byte, 2)
	binary.LittleEndian.PutUint16(messageLength, messageLengthUint)
	osdpPacket := &OSDPPacket{
		startOfMessage:    OSDPSOM,
		peripheralAddress: peripheralAddress, lsbLength: messageLength[0], msbLength: messageLength[1],
		msgCtrlInfo: msgControlInfo, securityBlockLength: 0x00, securityBlockType: 0x00, securityBlockData: nil,
		msgCode: byte(msgCode), msgData: msgData, msgAuthenticationCode: nil,
		lsbChecksum: 0x00, msbChecksum: 0x00, secure: false, useMAC: false,
	}

	osdpPacket.calculateCRC()
	return osdpPacket, nil
}

func (osdpPacket *OSDPPacket) calculateCRC() {
	osdpPacketBytes := osdpPacket.ToBytes()
	packetBytesSizeWithoutChecksum := len(osdpPacketBytes) - 2
	crc16Table := crc16.MakeTable(crc16.CRC16_AUG_CCITT)
	checksumUint := crc16.Checksum(osdpPacketBytes[:packetBytesSizeWithoutChecksum], crc16Table)
	checksum := make([]byte, 2)
	binary.LittleEndian.PutUint16(checksum, checksumUint)
	osdpPacket.lsbChecksum = checksum[0]
	osdpPacket.msbChecksum = checksum[1]
}

func (osdpPacket *OSDPPacket) ToBytes() []byte {
	var packetBytes []byte = []byte{
		osdpPacket.startOfMessage, osdpPacket.peripheralAddress,
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
	// Check that start of message follows OSDP spec
	currentIndex := 0
	startOfMessage := OSDPSOM
	for i := range payload {
		if payload[i] == startOfMessage {
			payload = payload[i:]
			break
		}
	}

	// Check that payload meets minimum OSDP spec size
	var payloadLength uint16 = uint16(len(payload))
	if payloadLength < minimumPacketLengthUnsecure {
		return nil, PacketIncompleteError
	}

	if payload[currentIndex] != startOfMessage {
		return nil, InvalidSOMError
	}

	currentIndex++
	// Check that the peripheral Address is in range
	peripheralAddress := payload[currentIndex]
	if (peripheralAddress&maxPeripheralAddress) < minPeripheralAddress || (peripheralAddress&maxPeripheralAddress) > maxPeripheralAddress {
		return nil, AddressOutOfRangeError
	}

	// Parse the message length
	currentIndex++
	messageLength := uint16(payload[currentIndex]) + uint16(payload[currentIndex+1])<<8
	bytesRemaining := messageLength - minimumPacketLengthUnsecure
	if len(payload) < int(messageLength) {
		return nil, PacketIncompleteError
	}

	// Check the message control info. TODO: Check for MAC
	currentIndex += 2
	msgControlInfo := payload[currentIndex]
	integrityCheck := false
	if (msgControlInfo & msgControlChecksumMask) == msgControlChecksumMask {
		integrityCheck = true
	}
	sequenceNumber := msgControlInfo & 0x03

	secure := (msgControlInfo & msgControlSecureMask) == msgControlSecureMask
	secureBlockType := byte(0x00)
	secureBlockLength := byte(0x00)
	secureBlockData := []byte{}
	// TODO check the security block length if secure
	if secure {
		currentIndex += 1
		secureBlockLength = payload[currentIndex]
		bytesRemaining -= uint16(secureBlockLength)
		currentIndex += 1
		secureBlockType = payload[currentIndex]
		secureBlockLength -= 0x02
		if len(payload) <= currentIndex+int(secureBlockLength) {
			return nil, PacketIncompleteError
		} else if secureBlockLength > 0x00 {
			currentIndex += 1
			secureBlockData = payload[currentIndex : currentIndex+int(secureBlockLength)]
		}
	}

	currentIndex++
	if len(payload) <= currentIndex {
		return nil, PacketIncompleteError
	}
	// Check the message code
	msgCode := payload[currentIndex]
	currentIndex++
	if len(payload) < (currentIndex + int(bytesRemaining)) {
		return nil, PacketIncompleteError
	}
	// TODO: if MAC then subtract 4 from bytes remaining to get length of msgData

	msgData := payload[currentIndex : currentIndex+int(bytesRemaining)]
	var MAC []byte
	if secureBlockType > SCS_14 {
		if len(msgData) < 4 {
			// MAC should be present but it's not
			return nil, PacketIncompleteError
		}
		MAC = msgData[len(msgData)-4:]
		msgData = msgData[:len(msgData)-4]
	}

	currentIndex += int(bytesRemaining)

	lsbChecksum := payload[currentIndex]
	currentIndex++
	msbChecksum := payload[currentIndex]

	if secure == false {
		osdpPacket, err := NewPacket(OSDPCode(msgCode), peripheralAddress, msgData, sequenceNumber, integrityCheck)
		if err != nil {
			return nil, err
		}

		if lsbChecksum != osdpPacket.lsbChecksum || msbChecksum != osdpPacket.msbChecksum {
			return nil, ChecksumFailedError
		}
		return osdpPacket, err
	}
	osdpPacket, err := NewSecurePacket(OSDPCode(msgCode), peripheralAddress, msgData, secureBlockType, secureBlockData, sequenceNumber, integrityCheck)
	if err != nil {
		return nil, err
	}
	osdpPacket.msgAuthenticationCode = MAC
	osdpPacket.recalculateChecksum()
	if lsbChecksum != osdpPacket.lsbChecksum || msbChecksum != osdpPacket.msbChecksum {
		return nil, ChecksumFailedError
	}

	return osdpPacket, err
}

func (osdpPacket *OSDPPacket) GetPeripheralAddress() byte {
	return osdpPacket.peripheralAddress
}

func (osdpPacket *OSDPPacket) GetMessageCode() byte {
	return osdpPacket.msgCode
}

func (osdpPacket *OSDPPacket) GetMessageData() []byte {
	return osdpPacket.msgData
}

func (osdpPacket *OSDPPacket) IsSecure() bool {
	return osdpPacket.secure
}

func (osdpPacket *OSDPPacket) GetSecurityBlockType() byte {
	return osdpPacket.securityBlockType
}

func (osdpPacket *OSDPPacket) GetSecurityBlockData() []byte {
	return osdpPacket.securityBlockData
}

func (osdpPacket *OSDPPacket) GetSequenceNumber() byte {
	sequenceNumber := osdpPacket.msgCtrlInfo & 0x03
	return sequenceNumber
}

func (osdpPacket *OSDPPacket) recalculateChecksum() {
	osdpPacketBytes := osdpPacket.ToBytes()
	packetBytesSizeWithoutChecksum := len(osdpPacketBytes) - 2
	crc16Table := crc16.MakeTable(crc16.CRC16_AUG_CCITT)
	checksumUint := crc16.Checksum(osdpPacketBytes[:packetBytesSizeWithoutChecksum], crc16Table)
	checksum := make([]byte, 2)
	binary.LittleEndian.PutUint16(checksum, checksumUint)
	osdpPacket.lsbChecksum = checksum[0]
	osdpPacket.msbChecksum = checksum[1]
}
