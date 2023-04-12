package transport

import (
	flowmessage "github.com/cloudflare/goflow/v3/pb"
)

type transport struct {
	consumers []chan<- []*flowmessage.FlowMessage
}

func NewTransport(consumers []chan<- []*flowmessage.FlowMessage) *transport {
	return &transport{
		consumers: consumers,
	}
}

func (t *transport) Publish(messages []*flowmessage.FlowMessage) {
	// Distribute the same slice to all consumers.
	// Consumers are expected to only read from it.
	for _, ch := range t.consumers {
		ch <- messages
	}
}
