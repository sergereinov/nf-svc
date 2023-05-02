package config

import (
	"errors"
	"os"
	"path/filepath"
)

const (
	_DEFAULT_INI_PATH  = "./nf-svc.ini"
	_DEFAULT_PORT      = 2055
	_DEFAULT_TOP_COUNT = 100
	_DEFAULT_KEEP_DAYS = 30
	_DEFAULT_MAX_SIZE  = 10
	_DEFAULT_LOGS_DIR  = "./logs"
)

var (
	_DEFAULT_INTERVALS = []int{20, 60, 8 * 60, 24 * 60}
)

func Load(optPath ...string) (path string, cfg *Config, err error) {
	path = getIniPath(optPath...)

	// Load ini file
	file, errLoad := NewIniFile(path)

	// Prepare config
	cfg = &Config{
		port: file.Int("Settings", "Port", _DEFAULT_PORT),

		Summary: Summary{
			intervals: file.Ints("Summary", "Intervals", _DEFAULT_INTERVALS),
			topCount:  file.Int("Summary", "TopCount", _DEFAULT_TOP_COUNT),
		},

		Logs: Logs{
			keepDays:         file.Int("Logs", "KeepDays", _DEFAULT_KEEP_DAYS),
			maxFileSizeMB:    file.Int("Logs", "MaxFileSizeMB", _DEFAULT_MAX_SIZE),
			dir:              file.String("Logs", "Dir", _DEFAULT_LOGS_DIR),
			enableSummaryLog: file.Bool("Logs", "EnableSummaryLog", true),
			enableNetFlowLog: file.Bool("Logs", "EnableNetFlowLog", true),
		},
	}

	errSave := file.SaveTo(path)

	// Update cfg logs path to absolute
	if !filepath.IsAbs(cfg.Logs.dir) && filepath.IsAbs(path) {
		dir := filepath.Dir(path)
		cfg.Logs.dir = filepath.Clean(filepath.Join(dir, cfg.Logs.dir))
	}

	return path, cfg, errors.Join(errLoad, errSave)
}

func getIniPath(optPath ...string) string {
	// Return absolute path if given
	if len(optPath) > 0 {
		path := optPath[0]
		if !filepath.IsAbs(path) {
			execPath, _ := os.Executable()
			path = filepath.Clean(filepath.Join(filepath.Dir(execPath), path))
		}
		return path
	}

	// Build ini path from executable location
	path, err := os.Executable()
	if err != nil || len(path) == 0 {
		return _DEFAULT_INI_PATH
	}

	// Remove ext
	for i := len(path) - 1; i >= 0 && !os.IsPathSeparator(path[i]); i-- {
		if path[i] == '.' {
			path = path[:i]
			break
		}
	}

	return path + ".ini"
}
