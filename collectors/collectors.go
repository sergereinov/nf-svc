package collectors

import (
	flowmessage "github.com/cloudflare/goflow/v3/pb"
)

type CollectorsConfig interface {
	GetSummaryIntervals() []int
	GetTrackingClients() []string
}

type Loggers struct {
	Summary chan<- string
	Netflow chan<- string
}

// Create collectors that will aggregate, transform and log netflow messages
func NewCollectors(cfg CollectorsConfig, logs Loggers) []chan<- []*flowmessage.FlowMessage {

	// Create collectors that will aggregate summaries
	summaryIntervals := cfg.GetSummaryIntervals()
	consumers := make([]chan<- []*flowmessage.FlowMessage, 0, len(summaryIntervals)+1)
	for _, interval := range summaryIntervals {
		if interval > 0 {
			c := NewSummaryCollector(interval, logs.Summary, cfg.GetTrackingClients())
			consumers = append(consumers, c.GetMessagesChannel())
		}
	}

	// Create collector that will thranform netflow messages
	c := NewNetflowCollector(logs.Netflow)
	consumers = append(consumers, c.GetMessagesChannel())

	return consumers
}
