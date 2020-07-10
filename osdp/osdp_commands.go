package osdp

type OSDPCommand byte

const (
	CMD_POLL     OSDPCommand = 0x60
	CMD_ID       OSDPCommand = 0x61
	CMD_CAP      OSDPCommand = 0x62
	CMD_DIAG     OSDPCommand = 0x63
	CMD_LSTAT    OSDPCommand = 0x64
	CMD_ISTAT    OSDPCommand = 0x65
	CMD_OSTAT    OSDPCommand = 0x66
	CMD_RSTAT    OSDPCommand = 0x67
	CMD_OUT      OSDPCommand = 0x68
	CMD_LED      OSDPCommand = 0x69
	CMD_BUZ      OSDPCommand = 0x6A
	CMD_TEXT     OSDPCommand = 0x6B
	CMD_COMSET   OSDPCommand = 0x6E
	CMD_DATA     OSDPCommand = 0x6F
	CMD_PROMPT   OSDPCommand = 0x71
	CMD_BIOREAD  OSDPCommand = 0x73
	CMD_BIOMATCH OSDPCommand = 0x74
	CMD_KEYSET   OSDPCommand = 0x75
	CMD_CHLNG    OSDPCommand = 0x76
	CMD_SCRYPT   OSDPCommand = 0x77
	CMD_ABORT    OSDPCommand = 0x7A
	CMD_MAXREPLY OSDPCommand = 0x7B
	CMD_MFG      OSDPCommand = 0x80
)
