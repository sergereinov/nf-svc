//go:build !windows

package service

import (
	"context"
	"os/signal"
	runtimeDebug "runtime/debug"
	"syscall"
)

/**************
 	Linux:
	For service's operations please use systemd service manager

	Install:
	sudo systemclt enable nf-svc.service
	sudo systemctl start nf-svc.service

	Uninstall:
	sudo systemctl stop nf-svc.service
	sudo systemclt disable nf-svc.service

--- /lib/systemd/system/nf-svc.service sample file ---
[Unit]
Description=Standalone NetFlow collector
ConditionPathExists=/home/user/nf-svc/nf-svc
After=network.target

[Service]
Type=simple
User=user
Group=user
LimitNOFILE=1024

Restart=always
RestartSec=5

WorkingDirectory=/home/user/nf-svc
ExecStart=/home/user/nf-svc/nf-svc

[Install]
WantedBy=multi-user.target

*********************/

func (s Service) entryPoint(payload func(context.Context)) {
	defer func() {
		if x := recover(); x != nil {
			s.logger.Fatalf("panic: %v\n%v", x, string(runtimeDebug.Stack()))
		}
	}()

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer cancel()

	payload(ctx)
}
