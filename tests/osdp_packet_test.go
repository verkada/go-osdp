package main

import (
	"testing"

	"github.com/HowardDunn/go-osdp/osdp"
	"github.com/stretchr/testify/require"
)

func Sum(a, b int) int {

	return a + b
}

func TestSum(t *testing.T) {
	tables := []struct {
		x int
		y int
		n int
	}{
		{1, 1, 2},
		{1, 2, 3},
		{2, 2, 4},
		{5, 2, 7},
	}

	for _, table := range tables {
		total := Sum(table.x, table.y)
		if total != table.n {
			t.Errorf("Sum of (%d+%d) was incorrect, got: %d, want: %d.", table.x, table.y, total, table.n)
		}
	}
}

func TestPacketCreationACK(t *testing.T) {
	osdpPacket, err := osdp.NewPacket(osdp.REPLY_ACK, 0x00, []byte{}, true)
	if err != nil {
		t.Errorf("Unable to Create OSDP Packet")
	}

	correctMessage := []byte{0x53, 0x00, 0x08, 0x00, 0x04, 0x40, 0x89, 0x8E}
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
	osdpPacket, err := osdp.NewPacket(osdp.CMD_POLL, 0x65, card, true)
	if err != nil {
		t.Errorf("Unable to Create OSDP Packet")
	}

	correctMessage := []byte{0x53, 0x00, 0x22, 0x00, 0x04, 0x50, 0x30, 0x30, 0x30,
		0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30,
		0x31, 0x31, 0x31, 0x30, 0x30, 0x30, 0x31, 0x31, 0x30,
		0x31, 0x30, 0x31, 0x30, 0x31, 0xFE, 0x40}
	osdpPacketBytes := osdpPacket.ToBytes()
	require.Equal(t, correctMessage, osdpPacketBytes)
}