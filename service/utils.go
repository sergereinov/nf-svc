package service

import (
	"os"
	"path/filepath"
)

func ExecutableFilename() string {
	// Get executable instance path
	name, err := os.Executable()
	if err != nil {
		name = os.Args[0]
	}

	// Extract filename
	name = filepath.Base(name)

	// Remove ext
	for i := len(name) - 1; i >= 0 && !os.IsPathSeparator(name[i]); i-- {
		if name[i] == '.' {
			name = name[:i]
			break
		}
	}

	return name
}
