//go:build !windows

package service

import (
	"context"
	"os/signal"
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
ConditionPathExists=/home/nf-svc/nf-svc
After=network.target

[Service]
Type=simple
User=pi
Group=pi
LimitNOFILE=1024

Restart=always
RestartSec=5

WorkingDirectory=/home/nf-svc
ExecStart=/home/nf-svc/nf-svc

[Install]
WantedBy=multi-user.target

*********************/

func (s Service) entryPoint(payload func(context.Context)) {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer cancel()

	payload(ctx)
}
