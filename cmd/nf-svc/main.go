package main

import (
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
	// Load config
	pathIni, cfg, errIni := config.Load()

	// Create a general purpose logger and change the default logger
	log := loggers.NewCommonLogger(&cfg.Logs)

	// Report instance status
	execPath, _ := os.Executable()
	log.Infof("Starting %v", execPath)
	if errIni != nil {
		log.Errorf("Load %s: %v", pathIni, errIni)
	} else {
		log.Infof("Load %s", pathIni)
	}
	log.Infof("Config: %+v", cfg)

	// Create and run log-writers goroutines
	dumpSummary := loggers.NewSummaryWriter(&cfg.Logs)
	dumpNetflow := loggers.NewNetflowWriter(&cfg.Logs)

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
	err := s.FlowRoutine(workers, addr, cfg.Port, reuse)
	if err != nil {
		log.Fatalf("Fatal error: could not listen to UDP (%v)", err)
	}
}
