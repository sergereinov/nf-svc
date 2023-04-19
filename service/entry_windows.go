//go:build windows

package service

import (
	"context"
	"fmt"
	"os"
	"strings"

	"golang.org/x/sys/windows/svc"
)

func (s Service) entryPoint(payload func(context.Context)) {

	if isInService, err := svc.IsWindowsService(); err != nil {
		s.Logger.Printf("failed to determine if we are running in service: %v", err)
		return
	} else if isInService {
		s.run(payload)
		return
	}

	if len(os.Args) < 2 {
		s.printUsage()
		return
	}

	cmd := strings.ToLower(os.Args[1])
	switch cmd {
	case "/d":
		// Debug with console
		s.runWithConsole(payload)
		return
	case "/i":
		// Install service
		if err := s.install(); err != nil {
			s.Logger.Printf("failed to install %s v%s: %v\n", s.Name, s.Version, err)
		}
	case "/u":
		// Uninstall service
		if err := s.uninstall(); err != nil {
			s.Logger.Printf("failed to uninstall %s v%s: %v\n", s.Name, s.Version, err)
		}
	default:
		s.printUsage()
	}
}

func (s Service) printUsage() {
	fmt.Fprintf(os.Stdout, "%s, v%s\n"+
		"Usage: %s </i>|</u>|</d>\n"+
		"\t/i - install service.\n"+
		"\t/u - uninstall service.\n"+
		"\t/d - debug with console.\n",
		s.Name, s.Version, os.Args[0])
}
