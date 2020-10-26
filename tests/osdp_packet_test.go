package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	osdp "github.com/verkada/go-osdp"
)

func TestPacketCreationACK(t *testing.T) {
	osdpPacket, err := osdp.NewPacket(osdp.REPLY_ACK, 0x00, []byte{}, true)
	if err != nil {
		t.Errorf("Unable to Create OSDP Packet")
	}

	correctMessage := []byte{0x53, 0x00, 0x08, 0x00, 0x04, 0x40, 0x89, 0x8E}
	osdpPacketBytes := osdpPacket.ToBytes()
	require.Equal(t, correctMessage, osdpPacketBytes)
}

func TestPacketCreationNAK(t *testing.T) {
	osdpPacket, err := osdp.NewPacket(osdp.REPLY_NAK, 0x00, []byte{0x01}, true)
	if err != nil {
		t.Errorf("Unable to Create OSDP Packet")
	}

	correctMessage := []byte{0x53, 0x00, 0x09, 0x00, 0x04, 0x41, 0x01, 0x07, 0x70}
	osdpPacketBytes := osdpPacket.ToBytes()
	require.Equal(t, correctMessage, osdpPacketBytes)
}

func TestPacketCreationPOLL(t *testing.T) {
	osdpPacket, err := osdp.NewPacket(osdp.CMD_POLL, 0x65, []byte{}, true)
	if err != nil {
		t.Errorf("Unable to Create OSDP Packet")
	}

	correctMessage := []byte{0x53, 0x65, 0x08, 0x00, 0x04, 0x60, 0x60, 0x90}
	osdpPacketBytes := osdpPacket.ToBytes()
	require.Equal(t, correctMessage, osdpPacketBytes)
}

func TestPacketCreationCardScan(t *testing.T) {
	card := []byte("00000000010011100011010101")
	osdpPacket, err := osdp.NewPacket(osdp.REPLY_RAW, 0x00, card, true)
	if err != nil {
		t.Errorf("Unable to Create OSDP Packet")
	}

	correctMessage := []byte{
		0x53, 0x00, 0x22, 0x00, 0x04, 0x50, 0x30, 0x30, 0x30,
		0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30,
		0x31, 0x31, 0x31, 0x30, 0x30, 0x30, 0x31, 0x31, 0x30,
		0x31, 0x30, 0x31, 0x30, 0x31, 0xFE, 0x40,
	}
	osdpPacketBytes := osdpPacket.ToBytes()
	require.Equal(t, correctMessage, osdpPacketBytes)
}

func TestPacketDecodeACK(t *testing.T) {
	msgToDecode := []byte{0x53, 0x00, 0x08, 0x00, 0x04, 0x40, 0x89, 0x8E}
	decodedPacket, err := osdp.NewPacketFromBytes(msgToDecode)
	if err != nil {
		t.Errorf("Unable to Decode OSDP Packet: %v", err.Error())
	}

	correctPacket, err := osdp.NewPacket(osdp.REPLY_ACK, 0x00, []byte{}, true)
	if err != nil {
		t.Errorf("Unable to Create OSDP Packet: %v", err.Error())
	}

	require.Equal(t, correctPacket, decodedPacket)
}

func TestPacketDecodeCardScan(t *testing.T) {
	msgToDecode := []byte{
		0x53, 0x00, 0x22, 0x00, 0x04, 0x50, 0x30, 0x30, 0x30,
		0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30,
		0x31, 0x31, 0x31, 0x30, 0x30, 0x30, 0x31, 0x31, 0x30,
		0x31, 0x30, 0x31, 0x30, 0x31, 0xFE, 0x40,
	}
	decodedPacket, err := osdp.NewPacketFromBytes(msgToDecode)
	if err != nil {
		t.Errorf("Unable to Decode OSDP Packet: %v", err.Error())
	}
	card := []byte("00000000010011100011010101")
	correctPacket, err := osdp.NewPacket(osdp.REPLY_RAW, 0x00, card, true)
	if err != nil {
		t.Errorf("Unable to Create OSDP Packet: %v", err)
	}

	require.Equal(t, correctPacket, decodedPacket)
}

func TestMessageReceive(t *testing.T) {
	transceiver := &MockTransceiver{timesCalled: 0}
	messenger := osdp.NewOSDPMessenger(transceiver)

	correctMessage := &osdp.OSDPMessage{MessageCode: 0x40, PeripheralAddress: 0x00, MessageData: []byte{}}
	message, err := messenger.ReceiveResponse(1 * time.Second)
	if err != nil {
		t.Errorf("Error while Receiving Message response: %v", err.Error())
	}
	require.Equal(t, correctMessage, message)
}

func TestMessageReceiveTimeout(t *testing.T) {
	transceiver := &SlowTransceiver{}
	messenger := osdp.NewOSDPMessenger(transceiver)
	_, err := messenger.ReceiveResponse(200 * time.Millisecond)
	require.Equal(t, osdp.OSDPReceiveTimeoutError, err)
	correctMessage := &osdp.OSDPMessage{MessageCode: 0x40, PeripheralAddress: 0x00, MessageData: []byte{}}
	message, err := messenger.ReceiveResponse(400 * time.Millisecond)
	if err != nil {
		t.Errorf("Error while Receiving Message response: %v", err.Error())
	}
	require.Equal(t, correctMessage, message)
}
