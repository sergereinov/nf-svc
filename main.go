package main

import (
	"fmt"
	"sort"

	flowmessage "github.com/cloudflare/goflow/v3/pb"
	"github.com/cloudflare/goflow/v3/utils"
	log "github.com/sirupsen/logrus"
)

const (
	workers = 1
	addr    = ""
	port    = 2055
	reuse   = false
)

var summary Summary[GroupMsg]

var trackingClients = []string{"192.168.255.1"}

type transport struct{}

func (t transport) Publish(messages []*flowmessage.FlowMessage) {
	for _, m := range messages {
		summary.Add(Msg{FlowMessage: m, trackingClients: trackingClients})
	}

	fmt.Printf("*** summary dump ***\n")
	for partition, groups := range summary.Dump() {
		fmt.Printf("%v\n", partition)

		type row struct {
			group string
			value GroupMsg
		}
		rows := make([]row, 0, len(groups))
		for group, data := range groups {
			rows = append(rows, row{group: group, value: data})
		}

		sort.Slice(rows, func(a, b int) bool {
			return rows[a].value.Bytes > rows[b].value.Bytes //reverse order, from big to small
		})

		for _, r := range rows {
			fmt.Printf("  %v = %+v\n", r.group, r.value)
		}
	}
}

func main() {
	log.Info("Starting")

	_ = &utils.DefaultLogTransport{}

	s := &utils.StateNetFlow{
		//Transport: &utils.DefaultLogTransport{},
		Transport: &transport{},
		Logger:    log.StandardLogger(),
	}

	err := s.FlowRoutine(workers, addr, port, reuse)
	if err != nil {
		log.Fatalf("Fatal error: could not listen to UDP (%v)", err)
	}
}
