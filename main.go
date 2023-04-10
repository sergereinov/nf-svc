package main

import (
	"github.com/cloudflare/goflow/v3/utils"
	log "github.com/sirupsen/logrus"
)

const (
	workers = 1
	addr    = ""
	port    = 2055
	reuse   = false
)

func main() {
	log.Info("Starting")

	s := &utils.StateNetFlow{
		Transport: &utils.DefaultLogTransport{},
		Logger:    log.StandardLogger(),
	}

	err := s.FlowRoutine(workers, addr, port, reuse)
	if err != nil {
		log.Fatalf("Fatal error: could not listen to UDP (%v)", err)
	}
}
