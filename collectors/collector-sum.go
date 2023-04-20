package collectors

import (
	"context"
	"fmt"
	"runtime/debug"
	"sort"
	"strings"
	"sync"
	"time"

	flowmessage "github.com/cloudflare/goflow/v3/pb"
	"github.com/sergereinov/nf-svc/loggers"
)

type summaryCollector struct {
	interval        int
	messages        chan []*flowmessage.FlowMessage
	summary         Summary[GroupMsg]
	summaryTopCount int
	trackingClients map[string]struct{}
	reports         chan<- string
	logger          CollectorsLogger
}

type summaryCollectorConfig struct {
	interval        int
	summaryTopCount int
	trackingClients []string
	reports         chan<- string
	logger          CollectorsLogger
}

func NewSummaryCollector(ctx context.Context, wg *sync.WaitGroup, cfg summaryCollectorConfig) *summaryCollector {
	trackingMap := make(map[string]struct{})
	for _, c := range cfg.trackingClients {
		trackingMap[c] = struct{}{}
	}

	c := &summaryCollector{
		interval:        cfg.interval,
		messages:        make(chan []*flowmessage.FlowMessage),
		summaryTopCount: cfg.summaryTopCount,
		trackingClients: trackingMap,
		reports:         cfg.reports,
		logger:          cfg.logger,
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
			c.summary = Summary[GroupMsg]{}

		case messages := <-c.messages:
			for _, m := range messages {
				c.summary.Add(Msg{FlowMessage: m, trackingClients: c.trackingClients})
			}
		}
	}
}

func (c *summaryCollector) dumpSummary() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("*** Summary for every %d minutes ***%s", c.interval, loggers.LineBreak))

	for partition, groups := range c.summary.Dump() {

		if len(groups) <= c.summaryTopCount {
			sb.WriteString(fmt.Sprintf("%s%s", partition, loggers.LineBreak))
		} else {
			sb.WriteString(fmt.Sprintf("%s top %v of %v%s",
				partition,
				c.summaryTopCount,
				len(groups),
				loggers.LineBreak))
		}

		type row struct {
			group string
			value GroupMsg
		}

		rows := make([]row, 0, len(groups))
		for group, data := range groups {
			rows = append(rows, row{group: group, value: data})
		}

		//sort in reverse order, from big to small
		sort.Slice(rows, func(a, b int) bool {
			return rows[a].value.Bytes > rows[b].value.Bytes
		})

		//cut to top count
		if len(rows) > c.summaryTopCount {
			rows = rows[:c.summaryTopCount]
		}

		for _, r := range rows {
			sb.WriteString(fmt.Sprintf("  %v, %+v%s", r.group, r.value, loggers.LineBreak))
		}
	}

	sb.WriteString(loggers.LineBreak)

	return sb.String()
}
