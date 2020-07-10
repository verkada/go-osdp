package osdp_example

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/HowardDunn/go-osdp/osdp"
)

var osdpMessenger *osdp.OSDPMessenger

const (
	osdpMessageFrequencyMS time.Duration = 200
	osdpMessageTimeout     time.Duration = 200
	peripheralAddress      byte          = 0x00
)

func handleOSDPResponse(osdpMessage *osdp.OSDPMessage) error {

	return nil
}

func startCommunication(ctx context.Context, outgoingMessageChan chan *osdp.OSDPMessage) {

	ticker := time.NewTicker(osdpMessageFrequencyMS * time.Millisecond)

	executeOSDPCycle := func(outgoingMessage *osdp.OSDPMessage, timeout time.Duration) {
		osdpResponse, err := osdpMessenger.SendAndReceive(outgoingMessage, timeout)
		if err != nil {
			log.Fatal("Unable to Send and Receive")
		}
		handleOSDPResponse(osdpResponse)
	}

	for {

		select {
		case <-ctx.Done():
			return
		case outgoingMessage := <-outgoingMessageChan:
			executeOSDPCycle(outgoingMessage, osdpMessageTimeout)
		case <-ticker.C:
			osdpMessage := osdp.NewOSDPMessage(osdp.CMD_POLL, peripheralAddress, nil)
			executeOSDPCycle(osdpMessage, osdpMessageTimeout)
		}
	}
}

func Run() {

	osdpMessenger = osdp.NewOSDPMessenger(nil)
	var (
		wg          sync.WaitGroup
		ctx, cancel = context.WithCancel(context.Background())
	)
	wg.Add(1)
	outgoingMessages := make(chan *osdp.OSDPMessage, 1)
	go func() {
		startCommunication(ctx, outgoingMessages)
		wg.Done()
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	s := <-ch
	log.Print("signal", s, "shutting down")
	cancel()

	wg.Wait()
}
