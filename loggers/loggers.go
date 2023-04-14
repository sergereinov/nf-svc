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

	_MAX_BACKUP = 0     // disabled
	_COMPRESS   = false // disabled

	_DEFAULT_BASENAME = "nf-svc"
)

type LoggersConfig interface {
	GetKeepDays() int
	GetMaxFileSizeMB() int
	GetPath() string
}

func NewLoggers(cfg LoggersConfig) (logger *commonLogger, netflow *lumberjack.Logger, summary *lumberjack.Logger) {
	basename := baseExecutableName()
	logger = newCommonLogger(cfg, basename)
	netflow = newNetflowLogger(cfg, basename)
	summary = newSummaryLogger(cfg, basename)
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

func newCommonLogger(cfg LoggersConfig, basename string) *commonLogger {
	logPath := filepath.Join(cfg.GetPath(), basename+_COMMON_LOG)
	logger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    cfg.GetMaxFileSizeMB(),
		MaxBackups: _MAX_BACKUP,
		MaxAge:     cfg.GetKeepDays(),
		Compress:   _COMPRESS,
	}
	return &commonLogger{
		logger: logger,
		bw:     &BufferedWriter{},
	}
}

func newNetflowLogger(cfg LoggersConfig, basename string) *lumberjack.Logger {
	logPath := filepath.Join(cfg.GetPath(), basename+_NETFLOW_LOG)
	return &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    cfg.GetMaxFileSizeMB(),
		MaxBackups: _MAX_BACKUP,
		MaxAge:     cfg.GetKeepDays(),
		Compress:   _COMPRESS,
	}
}

func newSummaryLogger(cfg LoggersConfig, basename string) *lumberjack.Logger {
	logPath := filepath.Join(cfg.GetPath(), basename+_SUMMARY_LOG)
	return &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    cfg.GetMaxFileSizeMB(),
		MaxBackups: _MAX_BACKUP,
		MaxAge:     cfg.GetKeepDays(),
		Compress:   _COMPRESS,
	}
}
