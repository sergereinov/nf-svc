//go:build windows

package service

import (
	"context"
	"sync"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
)

// ref https://pkg.go.dev/golang.org/x/sys/windows/svc

type handler struct {
	payload func(context.Context)
	logger  Logger
}

func (h *handler) Execute(args []string, requests <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	// Report: begins the start procedure
	changes <- svc.Status{State: svc.StartPending}

	// Make signaling context
	ctx, cancelContext := context.WithCancel(context.Background())
	defer cancelContext()

	// Run the payload goroutine
	var wgPayload sync.WaitGroup
	payloadStopped := make(chan struct{})
	if h.payload != nil {
		wgPayload.Add(1)
		go func() {
			defer wgPayload.Done()
			defer close(payloadStopped)
			h.payload(ctx)
		}()
	}

	// Report: service is running and accepting some commands
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

loop:
	for {
		select {
		case <-payloadStopped:
			h.logger.Printf("Payload stopped unexpectedly")
			break loop
		case r := <-requests:
			switch r.Cmd {
			case svc.Interrogate:
				changes <- r.CurrentStatus
			case svc.Stop:
				//h.logger.Printf("Service manager requests svc.Stop")
				break loop
			case svc.Shutdown:
				//h.logger.Printf("Service manager requests svc.Shutdown")
				break loop
			default:
				h.logger.Printf("Unexpected control request: %+v", r)
			}
		}
	}

	// Report: begins the stop procedure
	changes <- svc.Status{State: svc.StopPending}

	// Stopping
	cancelContext()
	wgPayload.Wait()
	return
}

func (s Service) run(payload func(context.Context)) {
	h := &handler{
		payload: payload,
		logger:  s.Logger,
	}
	if err := svc.Run(s.Name, h); err != nil {
		s.Logger.Printf("Fail to run %s: %v", s.Name, err)
	}
}

func (s Service) runWithConsole(payload func(context.Context)) {
	s.Logger.Printf("Console debug mode, press Ctrl-C to exit...")
	h := &handler{
		payload: payload,
		logger:  s.Logger,
	}
	if err := debug.Run(s.Name, h); err != nil {
		s.Logger.Printf("Fail to run %s: %v", s.Name, err)
	}
}
