package loggers

import (
	"context"
	"sync"
)

const (
	_NETFLOW_LOG = "-netflow.log"
)

// Create and run netflow log-writer goroutine
func NewNetflowWriter(ctx context.Context, wg *sync.WaitGroup, cfg LoggersConfig) chan<- string {
	logger := newLogger(cfg, baseExecutableName()+_NETFLOW_LOG)
	input := make(chan string)
	NewLoggerWriter(ctx, wg, input, logger)
	return input
}
