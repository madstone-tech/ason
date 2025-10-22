package xdg

import (
	"os"
	"path/filepath"
)

// ConfigHome returns the XDG_CONFIG_HOME directory for ason
func ConfigHome() (string, error) {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configHome = filepath.Join(homeDir, ".config")
	}

	return filepath.Join(configHome, "ason"), nil
}

// DataHome returns the XDG_DATA_HOME directory for ason
func DataHome() (string, error) {
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dataHome = filepath.Join(homeDir, ".local", "share")
	}

	return filepath.Join(dataHome, "ason"), nil
}

// CacheHome returns the XDG_CACHE_HOME directory for ason
func CacheHome() (string, error) {
	cacheHome := os.Getenv("XDG_CACHE_HOME")
	if cacheHome == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		cacheHome = filepath.Join(homeDir, ".cache")
	}

	return filepath.Join(cacheHome, "ason"), nil
}
