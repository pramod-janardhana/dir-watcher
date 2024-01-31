package config

import (
	"errors"
	"os"
	"path/filepath"
)

func FileWatcherConfigPath() (string, error) {
	dir := os.ExpandEnv("${programdata}\\DirWatcher\\")
	filename := "config.filewatcher.windows.json"

	configPath := filepath.Join(dir, filename)
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		return "", err
	}

	return configPath, nil
}

func generalConfigPath() (string, error) {
	dir := os.ExpandEnv("${programdata}\\DirWatcher\\")
	filename := "config.general.windows.json"

	configPath := filepath.Join(dir, filename)
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		return "", err
	}

	return configPath, nil
}

func JonSchedulerConfigPath() (string, error) {
	dir := os.ExpandEnv("${programdata}\\DirWatcher\\")
	filename := "config.jobscheduler.windows.json"

	configPath := filepath.Join(dir, filename)
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		return "", err
	}

	return configPath, nil
}
