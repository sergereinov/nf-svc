package loggers

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

type commonLogger struct {
	logger *lumberjack.Logger
	bw     bufferedWriter
}

// Logger that conforms to the goflow-logger / logrus-logger interface
// src-ref: goflow/utils/utils.go
/*
type Logger interface {
	Printf(string, ...interface{})
	Errorf(string, ...interface{})
	Warnf(string, ...interface{})
	Warn(...interface{})
	Error(...interface{})
	Debug(...interface{})
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Fatalf(string, ...interface{})
}
*/

const (
	_ERROR = "[ERROR] "
	_WARN  = "[WARN] "
	_DEBUG = "[DEBUG] "
	_INFO  = "[INFO] "
	_FATAL = "[FATAL] "
)

func (c *commonLogger) printString(text string) {
	// log to logfile
	c.bw.WriteWithHeaderAndLineBreak(c.logger, text)
	// std log to stdout
	log.Print(text)
}

func (c *commonLogger) Printf(format string, args ...interface{}) {
	if c == nil {
		return
	}
	text := fmt.Sprintf(format, args...)
	c.printString(text)
}

func (c *commonLogger) Errorf(format string, args ...interface{}) {
	if c == nil {
		return
	}
	text := fmt.Sprintf(_ERROR+format, args...)
	c.printString(text)
}

func (c *commonLogger) Warnf(format string, args ...interface{}) {
	if c == nil {
		return
	}
	text := fmt.Sprintf(_WARN+format, args...)
	c.printString(text)
}

func (c *commonLogger) Warn(args ...interface{}) {
	if c == nil {
		return
	}
	text := _WARN + fmt.Sprint(args...)
	c.printString(text)
}

func (c *commonLogger) Error(args ...interface{}) {
	if c == nil {
		return
	}
	text := _ERROR + fmt.Sprint(args...)
	c.printString(text)
}

func (c *commonLogger) Debug(args ...interface{}) {
	if c == nil {
		return
	}
	text := _DEBUG + fmt.Sprint(args...)
	c.printString(text)
}

func (c *commonLogger) Debugf(format string, args ...interface{}) {
	if c == nil {
		return
	}
	text := fmt.Sprintf(_DEBUG+format, args...)
	c.printString(text)
}

func (c *commonLogger) Info(args ...interface{}) {
	if c == nil {
		return
	}
	text := _INFO + fmt.Sprint(args...)
	c.printString(text)
}

func (c *commonLogger) Infof(format string, args ...interface{}) {
	if c == nil {
		return
	}
	text := fmt.Sprintf(_INFO+format, args...)
	c.printString(text)
}

func (c *commonLogger) Fatalf(format string, args ...interface{}) {
	if c == nil {
		return
	}
	text := fmt.Sprintf(_FATAL+format, args...)
	c.printString(text)

	// like log.Fatal
	os.Exit(1)
}
