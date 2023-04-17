package loggers

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

// logger https://github.com/natefinch/lumberjack

// Common consts
const (
	_MAX_BACKUP = 0     // disabled
	_COMPRESS   = false // disabled

	_DEFAULT_BASENAME = "nf-svc"
)

// Loggers config interface for all type loggers
type LoggersConfig interface {
	GetKeepDays() int
	GetMaxFileSizeMB() int
	GetDir() string
}

// Get executable instance file name without file extension
func baseExecutableName() string {
	var basename string
	exe, err := os.Executable()
	if err == nil {
		basename = filepath.Base(exe)
	}

	//remove ext
	for i := len(basename) - 1; i >= 0 && !os.IsPathSeparator(basename[i]); i-- {
		if basename[i] == '.' {
			basename = basename[:i]
		}
	}

	//validate basename
	if basename == "." || strings.Contains(basename, string(filepath.Separator)) {
		basename = _DEFAULT_BASENAME
	}
	return basename
}

// Create Lumberjack logger with given filename and params
func newLogger(cfg LoggersConfig, filename string) *lumberjack.Logger {
	logPath := filepath.Join(cfg.GetDir(), filename)
	return &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    cfg.GetMaxFileSizeMB(),
		MaxBackups: _MAX_BACKUP,
		MaxAge:     cfg.GetKeepDays(),
		Compress:   _COMPRESS,
	}
}
