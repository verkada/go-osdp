package osdp

import "errors"

var (
	PacketIncompleteError       = errors.New("OSDP Packet Incomplete")
	InvalidSOMError             = errors.New("OSDP Packet Invalid SOM")
	AddressOutOfRangeError      = errors.New("Peripheral Address Out of Range")
	ChecksumFailedError         = errors.New("Checksum Failed Error")
	OSDPReceiveTimeoutError     = errors.New("OSDPReceiveTimeout")
	SecureBlockDataLengthError  = errors.New("Secure Block Data Length too large")
	InvalidSequenceNumber       = errors.New("Invalid Sequence Number")
	IncorrectRandomNumberLength = errors.New("Invalid Random Number byte array length")
	InvalidSecureBlockType      = errors.New("Invalid Secure Block Type")
)
