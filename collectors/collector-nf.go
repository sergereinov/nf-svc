package collectors

import (
	"context"
	"sync"

	flowmessage "github.com/cloudflare/goflow/v3/pb"
)

type netflowCollector struct {
	logger   chan<- string
	messages chan []*flowmessage.FlowMessage
}

func NewNetflowCollector(ctx context.Context, wg *sync.WaitGroup, logger chan<- string) *netflowCollector {
	c := &netflowCollector{
		logger:   logger,
		messages: make(chan []*flowmessage.FlowMessage),
	}

	wg.Add(1)
	go c.loop(ctx, wg)

	return c
}

func (c *netflowCollector) GetMessagesChannel() chan<- []*flowmessage.FlowMessage {
	return c.messages
}

func (c *netflowCollector) loop(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case messages := <-c.messages:
			for _, m := range messages {
				c.logger <- FormatFlowMessage(m)
			}
		}
	}
}
