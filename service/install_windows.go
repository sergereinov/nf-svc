// go:build windows

package service

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/sys/windows/svc/mgr"
)

// Example:
// https://github.com/golang/sys/blob/master/windows/svc/example/install.go

const (
	_RECOVERY_DELAY = 1000 * time.Millisecond
)

func (s Service) install() error {
	// Get instance path
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("fail to get executable file path: %w", err)
	}

	// Connect to Windows Service Manager
	manager, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("fail to connect to Windows Service Manager: %w", err)
	}
	defer manager.Disconnect()

	// Check service existence
	service, err := manager.OpenService(s.Name)
	if err == nil {
		service.Close()
		return fmt.Errorf("service %s already exists", s.Name)
	}

	// Register service
	cfg := mgr.Config{
		StartType:    mgr.StartAutomatic,
		ErrorControl: mgr.ErrorNormal,
		DisplayName:  s.Name,
		Description:  s.Description,
	}
	service, err = manager.CreateService(s.Name, exe, cfg)
	if err != nil {
		return fmt.Errorf("fail to create service %s: %w", s.Name, err)
	}
	defer service.Close()

	// Additional config
	ra := []mgr.RecoveryAction{
		mgr.RecoveryAction{
			Type:  mgr.ServiceRestart,
			Delay: _RECOVERY_DELAY,
		},
	}
	err = service.SetRecoveryActions(ra, 0)
	if err != nil {
		return fmt.Errorf("fail to set recovery actions: %w", err)
	}

	s.Logger.Printf("Service %s (v%s) installed successfully.\n", s.Name, s.Version)

	// Send command to run service via Windows Service Manager
	err = service.Start()
	if err != nil {
		s.Logger.Printf("Fail to start %s service on install: %v\n", s.Name, err)
	}

	return nil
}
