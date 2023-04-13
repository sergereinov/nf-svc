package loggers

import (
	"io"
	"time"
)

type bufferedWriter struct {
	buf []byte
}

func (l bufferedWriter) WriteWithHeaderAndLineBreak(w io.Writer, text string) {
	now := time.Now()

	l.buf = l.buf[:0]
	l.buf = append(l.buf, now.Format("2006-01-02 15:04:05.000")...)
	l.buf = append(l.buf, ' ')
	l.buf = append(l.buf, text...)

	if len(text) == 0 || text[len(text)-1] != '\n' {
		l.buf = append(l.buf, LineBreak...)
	}

	w.Write(l.buf)
}
