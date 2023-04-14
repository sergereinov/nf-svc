package loggers

import "io"

type loggerWriter struct {
	input  <-chan string
	output io.Writer
	bw     *BufferedWriter
}

func NewLoggerWriter(input <-chan string, output io.Writer) *loggerWriter {
	writer := &loggerWriter{
		input:  input,
		output: output,
		bw:     &BufferedWriter{},
	}
	go writer.loop()
	return writer
}

func (w *loggerWriter) loop() {
	for s := range w.input {
		w.bw.WriteWithHeaderAndLineBreak(w.output, s)
	}
}
