package main

import (
	"crypto/aes"
	"testing"

	"github.com/stretchr/testify/require"
	osdp "github.com/verkada/go-osdp"
)

var (
	randomNumberCP = []byte{0xB0, 0xB1, 0xB2, 0xB3, 0xB4, 0xB5, 0xB6, 0xB7}
	defaultSCBK    = []byte{0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F}
)

func TestOSDPCHLNGCreation(t *testing.T) {
	// Use Sequence number 0x01 and SB Data 0x00 to signify using defaultSCBK
	chlngPacket, err := osdp.NewSecurePacket(osdp.CMD_CHLNG, 0x00, randomNumberCP, osdp.SCS_11, []byte{0x00}, 0x01, true)
	if err != nil {
		t.Errorf("Unable to Create OSDP Packet")
	}
	correctCHLNGPacket := []byte{0x53, 0x00, 0x13, 0x00, 0x0D, 0x03, 0x11, 0x00, 0x76, 0xB0, 0xB1, 0xB2, 0xB3, 0xB4, 0xB5, 0xB6, 0xB7, 0x31, 0x77}
	chlngPacketBytes := chlngPacket.ToBytes()
	require.Equal(t, correctCHLNGPacket, chlngPacketBytes)
}

func TestOSDPSCRYPTCreation(t *testing.T) {
	// PD Response
	PDResponse := []byte{
		0x53, 0x80, 0x2B, 0x00, 0x0D, 0x03, 0x12, 0x00, 0x76, 0x00, 0x06, 0x8E, 0x00, 0x00, 0x00, 0x00,
		0x00, 0xA0, 0xA1, 0xA2, 0xA3, 0xA4, 0xA5, 0xA6, 0xA7, 0xFD, 0xE5, 0xD2, 0xF4, 0x28, 0xEC, 0x16,
		0x31, 0x24, 0x71, 0xEA, 0x3C, 0x02, 0xBD, 0x77, 0x96, 0xF8, 0x1E,
	}

	// Calculate S_ENC
	S_ENC_RAW := []byte{
		0x01, 0x82, randomNumberCP[0], randomNumberCP[1], randomNumberCP[2], randomNumberCP[3], randomNumberCP[4], randomNumberCP[5],
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	S_MAC1_RAW := []byte{
		0x01, 0x01, randomNumberCP[0], randomNumberCP[1], randomNumberCP[2], randomNumberCP[3], randomNumberCP[4], randomNumberCP[5],
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	S_MAC2_RAW := []byte{
		0x01, 0x02, randomNumberCP[0], randomNumberCP[1], randomNumberCP[2], randomNumberCP[3], randomNumberCP[4], randomNumberCP[5],
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	cipherSCBK, err := aes.NewCipher(defaultSCBK)
	if err != nil {
		t.Errorf("Unable to create new cipher")
	}
	S_ENC := make([]byte, 16)
	S_MAC1 := make([]byte, 16)
	S_MAC2 := make([]byte, 16)
	cipherSCBK.Encrypt(S_ENC, S_ENC_RAW)
	cipherSCBK.Encrypt(S_MAC1, S_MAC1_RAW)
	cipherSCBK.Encrypt(S_MAC2, S_MAC2_RAW)

	correct_S_ENC := []byte{0xbf, 0x8d, 0xc2, 0xa8, 0x32, 0x9a, 0xcb, 0x8c, 0x67, 0xc6, 0xd0, 0xcd, 0x9a, 0x45, 0x16, 0x82}
	correct_S_MAC1 := []byte{0x5e, 0x86, 0xc6, 0x76, 0x60, 0x3b, 0xde, 0xe2, 0xd8, 0xbe, 0xaf, 0xe1, 0x78, 0x63, 0x73, 0x32}
	correct_S_MAC2 := []byte{0x6f, 0xda, 0x86, 0xe8, 0x57, 0x77, 0x7e, 0x81, 0x13, 0x20, 0x35, 0x75, 0x82, 0x39, 0x17, 0x2e}
	require.Equal(t, correct_S_ENC, S_ENC)
	require.Equal(t, correct_S_MAC1, S_MAC1)
	require.Equal(t, correct_S_MAC2, S_MAC2)

	// First Verify that the PD replies with SCS_12
	require.Equal(t, int(osdp.SCS_12), int(PDResponse[6]))

	// Verify the PD's Cryptogram
	// PD's cryptogram = AES128(append(randomNumberCP, randomNumberPD), S-ENC )
	randomNumberPD := PDResponse[17:25]
	randomNumberCombined := append(randomNumberCP, randomNumberPD...)
	cipherS_ENC, err := aes.NewCipher(S_ENC)
	if err != nil {
		t.Errorf("Unable to create new cipher")
	}
	clientCryptogramGiven := PDResponse[25:41]
	clientCryptogramGenerated := make([]byte, 16)
	cipherS_ENC.Encrypt(clientCryptogramGenerated, randomNumberCombined)
	require.Equal(t, clientCryptogramGiven, clientCryptogramGenerated)

	// Generate Server Cryptogram
	//  CP's cryptogram = AES128(append(randomNumberCP, randomNumberPD), S-ENC )
	randomNumberCombined = append(randomNumberPD, randomNumberCP...)
	serverCryptogram := make([]byte, 16)
	cipherS_ENC.Encrypt(serverCryptogram, randomNumberCombined)
	correct_osdp_SCRYPT_Packet := []byte{
		0x53, 0x00, 0x1B, 0x00, 0x0E, 0x03, 0x13, 0x00, 0x77, 0x26, 0xD3, 0x35, 0x6E,
		0x07, 0x76, 0x2D, 0x26, 0x28, 0x01, 0xFC, 0x8E, 0x66, 0x65, 0xA8, 0x91, 0x40, 0xB4,
	}
	// Use Sequence number 0x02 and SB Data 0x00 to signify using defaultSCBK
	scryptPacket, err := osdp.NewSecurePacket(osdp.CMD_SCRYPT, 0x00, serverCryptogram, osdp.SCS_13, []byte{0x00}, 0x02, true)
	if err != nil {
		t.Errorf("Unable to Create OSDP Packet")
	}

	require.Equal(t, correct_osdp_SCRYPT_Packet, scryptPacket.ToBytes())
}
