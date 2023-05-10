package main

import (
	"context"
	"os"
	"runtime/debug"
	"sync"

	service "github.com/sergereinov/go-windows-service"

	"github.com/sergereinov/nf-svc/collectors"
	"github.com/sergereinov/nf-svc/config"
	"github.com/sergereinov/nf-svc/loggers"
	"github.com/sergereinov/nf-svc/transport"

	"github.com/cloudflare/goflow/v3/utils"
)

// https://pkg.go.dev/cmd/link
// go build -ldflags="-X main.Version=1.0.0" ./cmd/nf-svc

var Version = "0.2"
var Name = service.ExecutableFilename()
var Description = "Standalone NetFlow collector"

const (
	workers = 1
	addr    = ""
	reuse   = false
)

func main() {
	// Load config
	pathIni, cfg, errIni := config.Load()

	// Create a general purpose logger
	log := loggers.NewCommonLogger(&cfg.Logs)

	// Run service wrapper
	service.Service{
		Version:     Version,
		Name:        Name,
		Description: Description,
		Logger:      log,
	}.Proceed(func(ctx context.Context) {

		// Report instance status
		execPath, _ := os.Executable()
		log.Printf("Starting %s, v%s (%v)", Name, Version, execPath)
		if errIni != nil {
			log.Errorf("Ini-file %s: %v", pathIni, errIni)
		} else {
			log.Infof("Ini-file %s", pathIni)
		}
		log.Infof("Config: %+v", *cfg)

		var wg sync.WaitGroup

		// Create netflow consumers
		collectorsLoggers := collectors.Loggers{
			Common: log,
		}
		if cfg.Logs.EnableSummaryLog() {
			collectorsLoggers.Summary = loggers.NewSummaryWriter(ctx, &wg, &cfg.Logs)
		}
		if cfg.Logs.EnableNetFlowLog() {
			collectorsLoggers.Netflow = loggers.NewNetflowWriter(ctx, &wg, &cfg.Logs)
		}
		consumers := collectors.NewCollectors(ctx, &wg, &cfg.Summary, &cfg.Logs, collectorsLoggers)

		// We don't have methods to stop goflow goroutine.
		// So it will be stopped when exiting the main goroutine.
		go func() {
			defer func() {
				if x := recover(); x != nil {
					log.Fatalf("panic: %v\n%v", x, string(debug.Stack()))
				}
			}()

			// Create goflow compatible transport
			transport := transport.NewTransport(consumers)

			// Init goflow's StateNetFlow
			s := &utils.StateNetFlow{
				Transport: transport,
				Logger:    log,
			}

			// Run goflow's FlowRoutine
			err := s.FlowRoutine(workers, addr, cfg.Port(), reuse)
			if err != nil {
				log.Fatalf("Fatal error: could not listen to UDP (%v)", err)
			}
		}()

		<-ctx.Done() // just marking what signal we are waiting for to exit
		wg.Wait()

		log.Printf("Stopped %s, v%s (%v)", Name, Version, execPath)
	})
}
