package main

import (
	"log"
	"os"

	"github.com/sergereinov/nf-svc/collectors"
	"github.com/sergereinov/nf-svc/config"
	"github.com/sergereinov/nf-svc/loggers"
	"github.com/sergereinov/nf-svc/transport"

	flowmessage "github.com/cloudflare/goflow/v3/pb"
	"github.com/cloudflare/goflow/v3/utils"
)

const (
	workers = 1
	addr    = ""
	reuse   = false
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		// Log fatal error to default logger
		log.Fatalf("Error: %v", err)
	}

	// Create loggers and change the default logger
	log, netflowLogger, summaryLogger := loggers.NewLoggers(&cfg.Logs)

	execPath, _ := os.Executable()
	log.Infof("Starting %v", execPath)
	log.Infof("Config: %+v", cfg)

	// Create and run log-writers goroutines
	dumpSummary := make(chan string)
	loggers.NewLoggerWriter(dumpSummary, summaryLogger)
	dumpNetflow := make(chan string)
	loggers.NewLoggerWriter(dumpNetflow, netflowLogger)

	// Create collectors that will aggregate summaries
	consumers := make([]chan<- []*flowmessage.FlowMessage, 0, len(cfg.SummaryIntervals)+1)
	for _, interval := range cfg.SummaryIntervals {
		if interval > 0 {
			c := collectors.NewSummaryCollector(interval, dumpSummary, cfg.TrackingClients)
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
	err = s.FlowRoutine(workers, addr, cfg.Port, reuse)
	if err != nil {
		log.Fatalf("Fatal error: could not listen to UDP (%v)", err)
	}
}
