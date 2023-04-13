package main

import (
	"github.com/sergereinov/nf-svc/collectors"
	"github.com/sergereinov/nf-svc/loggers"
	"github.com/sergereinov/nf-svc/transport"

	flowmessage "github.com/cloudflare/goflow/v3/pb"
	"github.com/cloudflare/goflow/v3/utils"
)

const (
	workers = 1
	addr    = ""
	port    = 2055
	reuse   = false

	logsPath = "./logs"
)

var trackingClients = []string{"192.168.255.1"}
var dumpIntervals = []int{2, 10, 60}

func main() {
	// Create loggers
	log, netflowLogger, summaryLogger := loggers.NewLoggers(logsPath)

	log.Info("Starting")

	// Create and run log-writers goroutines
	dumpSummary := make(chan string)
	loggers.NewLoggerWriter(dumpSummary, summaryLogger)
	dumpNetflow := make(chan string)
	loggers.NewLoggerWriter(dumpNetflow, netflowLogger)

	// Create collectors that will aggregate summaries
	consumers := make([]chan<- []*flowmessage.FlowMessage, 0, len(dumpIntervals)+1)
	for _, interval := range dumpIntervals {
		if interval > 0 {
			c := collectors.NewSummaryCollector(interval, dumpSummary, trackingClients)
			consumers = append(consumers, c.GetMessagesChannel())
		}
	}

	// Create collector that will thranform netflow messages
	c := collectors.NewNetflowCollector(dumpNetflow)
	consumers = append(consumers, c.GetMessagesChannel())

	// Create transport that will distributes messages to collectors
	transport := transport.NewTransport(consumers)

	// Init goflow's StateNetFlow
	s := &utils.StateNetFlow{
		Transport: transport,
		Logger:    log,
	}

	// Run goflow's FlowRoutine
	err := s.FlowRoutine(workers, addr, port, reuse)
	if err != nil {
		log.Fatalf("Fatal error: could not listen to UDP (%v)", err)
	}
}
