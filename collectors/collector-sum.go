package collectors

import (
	"fmt"
	"sort"
	"strings"
	"time"

	flowmessage "github.com/cloudflare/goflow/v3/pb"
	"github.com/sergereinov/nf-svc/loggers"
)

type summaryCollector struct {
	interval        int
	dump            chan<- string
	messages        chan []*flowmessage.FlowMessage
	summary         Summary[GroupMsg]
	trackingClients map[string]struct{}
}

func NewSummaryCollector(interval int, dump chan<- string, tracking []string) *summaryCollector {
	trackingMap := make(map[string]struct{})
	for _, c := range tracking {
		trackingMap[c] = struct{}{}
	}

	c := &summaryCollector{
		interval:        interval,
		dump:            dump,
		messages:        make(chan []*flowmessage.FlowMessage),
		trackingClients: trackingMap,
	}
	go c.loop()
	return c
}

func (c *summaryCollector) GetMessagesChannel() chan<- []*flowmessage.FlowMessage {
	return c.messages
}

func (c *summaryCollector) loop() {
	ticker := time.NewTicker(time.Duration(c.interval) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.dump <- c.dumpSummary()
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
		sb.WriteString(fmt.Sprintf("%s%s", partition, loggers.LineBreak))

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

		for _, r := range rows {
			sb.WriteString(fmt.Sprintf("  %v, %+v%s", r.group, r.value, loggers.LineBreak))
		}
	}

	sb.WriteString(loggers.LineBreak)

	return sb.String()
}