package collectors

import (
	"context"
	"runtime/debug"
	"sync"

	flowmessage "github.com/cloudflare/goflow/v3/pb"
)

type netflowCollector struct {
	messages chan []*flowmessage.FlowMessage
	reports  chan<- string
	logger   CollectorsLogger
}

type netflowCollectorConfig struct {
	reports chan<- string
	logger  CollectorsLogger
}

func NewNetflowCollector(ctx context.Context, wg *sync.WaitGroup, cfg netflowCollectorConfig) *netflowCollector {
	c := &netflowCollector{
		messages: make(chan []*flowmessage.FlowMessage),
		reports:  cfg.reports,
		logger:   cfg.logger,
	}

	wg.Add(1)
	go c.loop(ctx, wg)

	return c
}

func (c *netflowCollector) GetMessagesChannel() chan<- []*flowmessage.FlowMessage {
	return c.messages
}

func (c *netflowCollector) loop(ctx context.Context, wg *sync.WaitGroup) {
	defer func() {
		if x := recover(); x != nil {
			c.logger.Fatalf("panic: %v\n%v", x, string(debug.Stack()))
		}
	}()
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case messages := <-c.messages:
			for _, m := range messages {
				c.reports <- FormatFlowMessage(m)
			}
		}
	}
}
