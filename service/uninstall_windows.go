//go:build windows

package service

import (
	"fmt"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

func (s Service) uninstall() error {

	// Connect to Windows Service Manager
	manager, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("fail to connect to Windows Service Manager: %w", err)
	}
	defer manager.Disconnect()

	// Open service interface
	service, err := manager.OpenService(s.Name)
	if err != nil {
		return fmt.Errorf("service %s is not installed", s.Name)
	}
	defer service.Close()

	// Send command to stop service
	_, err = service.Control(svc.Stop)
	if err != nil {
		s.Logger.Printf("Fail to stop %s service on uninstall: %v\n", s.Name, err)
	}

	// Unregister service
	err = service.Delete()
	if err != nil {
		return fmt.Errorf("fail to delete %s service: %w", s.Name, err)
	}

	s.Logger.Printf("Service %s uninstalled successfully.\n", s.Name)
	return nil
}
