package loggers

import (
	"context"
	"io"
	"sync"
)

type loggerWriter struct {
	input  <-chan string
	output io.Writer
	bw     *BufferedWriter
}

func NewLoggerWriter(ctx context.Context, wg *sync.WaitGroup, input <-chan string, output io.Writer) *loggerWriter {
	writer := &loggerWriter{
		input:  input,
		output: output,
		bw:     &BufferedWriter{},
	}

	wg.Add(1)
	go writer.loop(ctx, wg)
	return writer
}

func (w *loggerWriter) loop(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case s := <-w.input:
			w.bw.WriteWithHeaderAndLineBreak(w.output, s)
		}
	}
}
