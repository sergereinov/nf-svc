package collectors

import (
	"context"
	"sync"

	flowmessage "github.com/cloudflare/goflow/v3/pb"
)

type SummaryConfig interface {
	Intervals() []int
	TopCount() int
}

type LoggersConfig interface {
	EnableSummaryLog() bool
	EnableNetFlowLog() bool
}

type CollectorsLogger interface {
	Printf(string, ...interface{})
	Warnf(string, ...interface{})
	Fatalf(string, ...interface{})
}

type Loggers struct {
	Common  CollectorsLogger
	Summary chan<- string
	Netflow chan<- string
}

// Create collectors that will aggregate, transform and log netflow messages
func NewCollectors(ctx context.Context, wg *sync.WaitGroup, sumcfg SummaryConfig, logcfg LoggersConfig, logs Loggers) []chan<- []*flowmessage.FlowMessage {

	var consumers []chan<- []*flowmessage.FlowMessage

	// Create collectors that will aggregate summaries
	if logcfg.EnableSummaryLog() && logs.Summary != nil {
		summaryIntervals := sumcfg.Intervals()
		//consumers := make([]chan<- []*flowmessage.FlowMessage, 0, len(summaryIntervals)+1)
		for _, interval := range summaryIntervals {
			if interval > 0 {
				logs.Common.Printf("Running collector routine: NetFlow Summary with interval %d minutes.", interval)
				c := NewSummaryCollector(ctx,
					wg,
					summaryCollectorConfig{
						interval: interval,
						topCount: sumcfg.TopCount(),
						reports:  logs.Summary,
						logger:   logs.Common,
					},
				)
				consumers = append(consumers, c.GetMessagesChannel())
			}
		}
	} else {
		logs.Common.Printf("Disabled collectors: NetFlow summary.")
	}

	// Create collector that will thranform netflow messages
	if logcfg.EnableNetFlowLog() && logs.Netflow != nil {
		logs.Common.Printf("Running collector routine: logging NetFlow packets.")
		c := NewNetflowCollector(ctx,
			wg,
			netflowCollectorConfig{
				reports: logs.Netflow,
				logger:  logs.Common,
			})
		consumers = append(consumers, c.GetMessagesChannel())
	} else {
		logs.Common.Printf("Disabled collector: logging NetFlow packets.")
	}

	if len(consumers) == 0 {
		logs.Common.Warnf("There is no enabled NetFlow loggers / no NetFlow collectors are running!")
	}

	return consumers
}
