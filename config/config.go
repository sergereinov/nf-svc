package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

type Config struct {
	IniPath          string
	Port             int
	TrackingClients  []string
	SummaryIntervals []int
	Logs             Logs
}

const (
	_DEFAULT_INI = "./nf-svc.ini"
)

var (
	_DEFAULT_INTERVALS = []int{20, 60, 8 * 60}
)

func Load() (*Config, error) {
	path := getIniPath()

	cfg, err := ini.Load(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	settings := cfg.Section("Settings")

	port := settings.Key("Port").MustInt(2055)

	clientsSlice := settings.Key("TrackingClients").Strings(",")

	summaryIntervals := settings.Key("SummaryIntervals").Ints(",")
	if len(summaryIntervals) == 0 {
		summaryIntervals = _DEFAULT_INTERVALS
	}

	logs := cfg.Section("Logs")
	logsKeepDays := logs.Key("KeepDays").MustInt(30)
	logsMaxFileSize := logs.Key("MaxFileSizeMB").MustInt(10)
	logsPath := logs.Key("Path").MustString("./logs")
	if !filepath.IsAbs(logsPath) && filepath.IsAbs(path) {
		dir := filepath.Dir(path)
		logsPath = filepath.Join(dir, logsPath)
		logsPath = filepath.Clean(logsPath)
	}

	return &Config{
			IniPath:          path,
			Port:             port,
			TrackingClients:  clientsSlice,
			SummaryIntervals: summaryIntervals,
			Logs: Logs{
				KeepDays:      logsKeepDays,
				MaxFileSizeMB: logsMaxFileSize,
				Path:          logsPath,
			},
		},
		nil
}

func getIniPath() string {
	path, err := os.Executable()
	if err != nil {
		return _DEFAULT_INI
	}

	//remove ext
	for i := len(path) - 1; i >= 0 && !os.IsPathSeparator(path[i]); i-- {
		if path[i] == '.' {
			path = path[:i]
		}
	}

	return path + ".ini"
}
