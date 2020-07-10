package osdp

type OSDPResponse byte

const (
	REPLY_ACK       OSDPResponse = 0x40
	REPLY_NAK       OSDPResponse = 0x41
	REPLY_PDID      OSDPResponse = 0x45
	REPLY_PDCAP     OSDPResponse = 0x46
	REPLY_LSTATR    OSDPResponse = 0x48
	REPLY_IASTR     OSDPResponse = 0x49
	REPLY_OSTATR    OSDPResponse = 0x4A
	REPLY_RSTATR    OSDPResponse = 0x4B
	REPLY_RAW       OSDPResponse = 0x50
	REPLY_FMT       OSDPResponse = 0x51
	REPLY_KEYPAD    OSDPResponse = 0x53
	REPLY_COM       OSDPResponse = 0x54
	REPLY_BIOREADR  OSDPResponse = 0x57
	REPLY_BIOMATCHR OSDPResponse = 0x58
	REPLY_CCRYPT    OSDPResponse = 0x76
	REPLY_MFGREP    OSDPResponse = 0x90
	REPLY_BUSY      OSDPResponse = 0x79
	REPLY_XRD       OSDPResponse = 0xB1
)
