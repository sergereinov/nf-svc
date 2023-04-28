package collectors

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	flowmessage "github.com/cloudflare/goflow/v3/pb"
	"github.com/sergereinov/nf-svc/loggers"
	"github.com/sergereinov/nf-svc/summary"
)

type summarizer interface {
	Add(fmsg *flowmessage.FlowMessage)
	Reset()
	Dump(maxCount int, lineBreak string) string
}

type summaryCollector struct {
	interval int
	messages chan []*flowmessage.FlowMessage
	summary  summarizer
	topCount int
	reports  chan<- string
	logger   CollectorsLogger
}

type summaryCollectorConfig struct {
	interval int
	topCount int
	reports  chan<- string
	logger   CollectorsLogger
}

func NewSummaryCollector(ctx context.Context, wg *sync.WaitGroup, cfg summaryCollectorConfig) *summaryCollector {
	c := &summaryCollector{
		interval: cfg.interval,
		messages: make(chan []*flowmessage.FlowMessage),
		summary:  summary.NewSummary(),
		topCount: cfg.topCount,
		reports:  cfg.reports,
		logger:   cfg.logger,
	}

	wg.Add(1)
	go c.loop(ctx, wg)

	return c
}

func (c *summaryCollector) GetMessagesChannel() chan<- []*flowmessage.FlowMessage {
	return c.messages
}

func (c *summaryCollector) loop(ctx context.Context, wg *sync.WaitGroup) {
	defer func() {
		if x := recover(); x != nil {
			c.logger.Fatalf("panic: %v\n%v", x, string(debug.Stack()))
		}
	}()
	defer wg.Done()

	ticker := time.NewTicker(time.Duration(c.interval) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			c.reports <- c.dumpSummary()
			c.summary.Reset()

		case messages := <-c.messages:
			for _, msg := range messages {
				c.summary.Add(msg)
			}
		}
	}
}

func (c *summaryCollector) dumpSummary() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("*** Summary for every %d minutes ***%s", c.interval, loggers.LineBreak))
	sb.WriteString(c.summary.Dump(c.topCount, loggers.LineBreak))
	return sb.String()
}
