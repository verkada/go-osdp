package osdp

import "errors"

var (
	PacketIncompleteError   = errors.New("OSDP Packet Incomplete")
	InvalidSOMError         = errors.New("OSDP Packet Invalid SOM")
	AddressOutOfRangeError  = errors.New("Peripheral Address Out of Range")
	ChecksumFailedError     = errors.New("Checksum Failed Error")
	OSDPReceiveTimeoutError = errors.New("OSDPReceiveTimeout")
)
