package loggers

import (
	"context"
	"sync"
)

const (
	_SUMMARY_LOG = "-summary.log"
)

// Create and run summary log-writer goroutine
func NewSummaryWriter(ctx context.Context, wg *sync.WaitGroup, cfg LoggersConfig) chan<- string {
	logger := newLogger(cfg, baseExecutableName()+_SUMMARY_LOG)
	input := make(chan string)
	NewLoggerWriter(ctx, wg, input, logger)
	return input
}
