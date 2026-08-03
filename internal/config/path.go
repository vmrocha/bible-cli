package config

import (
	"errors"
	"os"
	"path/filepath"
)

const (
	appDirectory = "bible-cli"
	fileName     = "config.json"
)

// DefaultPath resolves the configuration file from the process environment.
func DefaultPath() (string, error) {
	configHome := os.Getenv("BIBLE_CLI_CONFIG_HOME")
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	home := ""
	if configHome == "" && xdgConfigHome == "" {
		var err error
		home, err = os.UserHomeDir()
		if err != nil {
			return "", errors.New("resolve configuration path: user home directory is unavailable")
		}
	}
	return ResolvePath(configHome, xdgConfigHome, home)
}

// ResolvePath applies the cross-platform Bible CLI configuration order.
func ResolvePath(configHome, xdgConfigHome, userHome string) (string, error) {
	switch {
	case configHome != "":
		if !filepath.IsAbs(configHome) {
			return "", errors.New("BIBLE_CLI_CONFIG_HOME must be an absolute path")
		}
		return filepath.Join(configHome, fileName), nil
	case xdgConfigHome != "":
		if !filepath.IsAbs(xdgConfigHome) {
			return "", errors.New("XDG_CONFIG_HOME must be an absolute path")
		}
		return filepath.Join(xdgConfigHome, appDirectory, fileName), nil
	case userHome == "":
		return "", errors.New("resolve configuration path: user home directory is unavailable")
	case !filepath.IsAbs(userHome):
		return "", errors.New("resolve configuration path: user home directory must be absolute")
	default:
		return filepath.Join(userHome, ".config", appDirectory, fileName), nil
	}
}
