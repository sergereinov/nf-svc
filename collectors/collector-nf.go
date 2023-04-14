package collectors

import (
	flowmessage "github.com/cloudflare/goflow/v3/pb"
)

type netflowCollector struct {
	dump     chan<- string
	messages chan []*flowmessage.FlowMessage
}

func NewNetflowCollector(dump chan<- string) *netflowCollector {
	c := &netflowCollector{
		dump:     dump,
		messages: make(chan []*flowmessage.FlowMessage),
	}
	go c.loop()
	return c
}

func (c *netflowCollector) GetMessagesChannel() chan<- []*flowmessage.FlowMessage {
	return c.messages
}

func (c *netflowCollector) loop() {
	for messages := range c.messages {
		for _, m := range messages {
			c.dump <- FormatFlowMessage(m)
		}
	}
}
