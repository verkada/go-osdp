package osdp_application

import (
	"context"
	"time"

	osdp "github.com/verkada/go-osdp"
)

var osdpMessenger *osdp.OSDPMessenger

const (
	osdpMessageFrequencyMS time.Duration = 200
	osdpMessageTimeout     time.Duration = 200
)

type (
	OSDPMessageHandler func(osdpMessage *osdp.OSDPMessage)
	ErrorHandler       func(err error)
)

func StartCommunication(ctx context.Context, transceiver osdp.OSDPTransceiver, osdpHandler OSDPMessageHandler, outgoingMessageChan chan *osdp.OSDPMessage, errorHandler ErrorHandler, peripheralAddress byte) {
	ticker := time.NewTicker(osdpMessageFrequencyMS * time.Millisecond)

	osdpMessenger = osdp.NewOSDPMessenger(transceiver)
	executeOSDPCycle := func(outgoingMessage *osdp.OSDPMessage, writeTimeout time.Duration, readTimeout time.Duration) {
		osdpResponse, err := osdpMessenger.SendAndReceive(outgoingMessage, writeTimeout, readTimeout)
		if err != nil {
			errorHandler(err)
			return
		}
		osdpHandler(osdpResponse)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case outgoingMessage := <-outgoingMessageChan:
			executeOSDPCycle(outgoingMessage, osdpMessageTimeout, osdpMessageTimeout)
		case <-ticker.C:
			osdpMessage := osdp.NewOSDPMessage(osdp.CMD_POLL, peripheralAddress, nil)
			executeOSDPCycle(osdpMessage, osdpMessageTimeout, osdpMessageTimeout)
		}
	}
}
