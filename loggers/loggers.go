package loggers

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

// logger https://github.com/natefinch/lumberjack

const (
	_COMMON_LOG  = ".log"
	_NETFLOW_LOG = "-netflow.log"
	_SUMMARY_LOG = "-summary.log"

	_MAX_SIZE   = 10    // megabytes
	_MAX_BACKUP = 0     // disabled
	_MAX_AGE    = 30    // days
	_COMPRESS   = false // disabled

	_DEFAULT_BASENAME = "nf-svc"
)

func NewLoggers(path string) (logger *commonLogger, netflow *lumberjack.Logger, summary *lumberjack.Logger) {
	basename := baseExecutableName()
	logger = newCommonLogger(path, basename)
	netflow = newNetflowLogger(path, basename)
	summary = newSummaryLogger(path, basename)
	return
}

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

func newCommonLogger(path, basename string) *commonLogger {
	logPath := filepath.Join(path, basename+_COMMON_LOG)
	logger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    _MAX_SIZE,
		MaxBackups: _MAX_BACKUP,
		MaxAge:     _MAX_AGE,
		Compress:   _COMPRESS,
	}
	return &commonLogger{
		logger: logger,
	}
}

func newNetflowLogger(path, basename string) *lumberjack.Logger {
	logPath := filepath.Join(path, basename+_NETFLOW_LOG)
	return &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    _MAX_SIZE,
		MaxBackups: _MAX_BACKUP,
		MaxAge:     _MAX_AGE,
		Compress:   _COMPRESS,
	}
}

func newSummaryLogger(path, basename string) *lumberjack.Logger {
	logPath := filepath.Join(path, basename+_SUMMARY_LOG)
	return &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    _MAX_SIZE,
		MaxBackups: _MAX_BACKUP,
		MaxAge:     _MAX_AGE,
		Compress:   _COMPRESS,
	}
}
