package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sergereinov/nf-svc/collectors"
	"github.com/sergereinov/nf-svc/config"
	"github.com/sergereinov/nf-svc/loggers"
	"github.com/sergereinov/nf-svc/transport"

	"github.com/cloudflare/goflow/v3/utils"
)

const (
	workers = 1
	addr    = ""
	reuse   = false
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer cancel()

	var wg sync.WaitGroup

	// Load config
	pathIni, cfg, errIni := config.Load()

	// Create a general purpose logger
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

	// Create netflow consumers
	collectorsLoggers := collectors.Loggers{
		Summary: loggers.NewSummaryWriter(ctx, &wg, &cfg.Logs),
		Netflow: loggers.NewNetflowWriter(ctx, &wg, &cfg.Logs),
	}
	consumers := collectors.NewCollectors(ctx, &wg, cfg, collectorsLoggers)

	// We don't have methods to stop goflow goroutine.
	// So it will be stopped when exiting the main goroutine.
	go func() {
		// Create goflow compatible transport
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
	}()

	<-ctx.Done() // just marking what signal we are waiting for to exit
	wg.Wait()

	log.Infof("Stopped %v", execPath)
}
