package main

import (
	"fmt"

	"github.com/sergereinov/nf-svc/collectors"
	"github.com/sergereinov/nf-svc/transport"

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

var trackingClients = []string{"192.168.255.1"}
var dumpIntervals = []int{2, 10, 60}

func main() {
	log.Info("Starting")

	dumpStrings := make(chan string)
	go func() {
		for s := range dumpStrings {
			fmt.Print(s)
		}
	}()

	consumers := make([]chan<- []*flowmessage.FlowMessage, 0, len(dumpIntervals))
	for _, v := range dumpIntervals {
		c := collectors.NewSummaryCollector(v, dumpStrings, trackingClients)
		consumers = append(consumers, c.GetMessagesChannel())
	}

	transport := transport.NewTransport(consumers)

	s := &utils.StateNetFlow{
		Transport: transport,
		Logger:    log.StandardLogger(),
	}

	err := s.FlowRoutine(workers, addr, port, reuse)
	if err != nil {
		log.Fatalf("Fatal error: could not listen to UDP (%v)", err)
	}
}
