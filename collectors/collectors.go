package collectors

import (
	"context"
	"sync"

	flowmessage "github.com/cloudflare/goflow/v3/pb"
)

type CollectorsConfig interface {
	GetSummaryIntervals() []int
	GetSummaryTopCount() int
	GetTrackingClients() []string
}

type Loggers struct {
	Summary chan<- string
	Netflow chan<- string
}

// Create collectors that will aggregate, transform and log netflow messages
func NewCollectors(ctx context.Context, wg *sync.WaitGroup, cfg CollectorsConfig, logs Loggers) []chan<- []*flowmessage.FlowMessage {

	// Create collectors that will aggregate summaries
	summaryIntervals := cfg.GetSummaryIntervals()
	consumers := make([]chan<- []*flowmessage.FlowMessage, 0, len(summaryIntervals)+1)
	for _, interval := range summaryIntervals {
		if interval > 0 {
			c := NewSummaryCollector(ctx,
				wg,
				summaryCollectorConfig{
					interval:        interval,
					summaryTopCount: cfg.GetSummaryTopCount(),
					trackingClients: cfg.GetTrackingClients(),
				},
				logs.Summary,
			)
			consumers = append(consumers, c.GetMessagesChannel())
		}
	}

	// Create collector that will thranform netflow messages
	c := NewNetflowCollector(ctx, wg, logs.Netflow)
	consumers = append(consumers, c.GetMessagesChannel())

	return consumers
}
