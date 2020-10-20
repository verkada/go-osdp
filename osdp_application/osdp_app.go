package osdp_application

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/verkada/go-osdp/osdp"
)

var osdpMessenger *osdp.OSDPMessenger

const (
	osdpMessageFrequencyMS time.Duration = 200
	osdpMessageTimeout     time.Duration = 200
	peripheralAddress      byte          = 0x00
)

type OSDPMessageHandler func(osdpMessage *osdp.OSDPMessage)

func StartCommunication(ctx context.Context, transceiver osdp.OSDPTransceiver, osdpHandler OSDPMessageHandler, outgoingMessageChan chan *osdp.OSDPMessage) {
	ticker := time.NewTicker(osdpMessageFrequencyMS * time.Millisecond)

	osdpMessenger = osdp.NewOSDPMessenger(transceiver)
	executeOSDPCycle := func(outgoingMessage *osdp.OSDPMessage, writeTimeout time.Duration, readTimeout time.Duration) {
		osdpResponse, err := osdpMessenger.SendAndReceive(outgoingMessage, writeTimeout, readTimeout)
		if err != nil {
			log.Error("Unable to Send and Receive: ", err)
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
