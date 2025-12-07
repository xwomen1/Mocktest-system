package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// the configuration directory
func GetConfigDir() (string, error) {
	if _, err := os.Stat("configs"); err == nil {
		return "configs", nil
	}

	exe, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exe)
		configDir := filepath.Join(exeDir, "configs")
		if _, err := os.Stat(configDir); err == nil {
			return configDir, nil
		}
	}

	home, err := os.UserHomeDir()
	if err == nil {
		configDir := filepath.Join(home, ".upm", "configs")
		if _, err := os.Stat(configDir); err == nil {
			return configDir, nil
		}
	}

	return ".", nil
}

func GetConfigFile(env string) (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	var configFile string
	switch strings.ToLower(env) {
	case "development", "dev":
		configFile = filepath.Join(configDir, "dev", "config.yaml")
	case "production", "prod":
		configFile = filepath.Join(configDir, "prod", "config.yaml")
	case "test":
		configFile = filepath.Join(configDir, "test", "config.yaml")
	default:
		configFile = filepath.Join(configDir, "config.yaml")
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {

		configFile = filepath.Join(configDir, "config.yaml")
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			return "", fmt.Errorf("config file not found: %s", configFile)
		}
	}

	return configFile, nil
}

func LoadConfig(env string) (*Config, error) {
	loader := NewLoader()

	if env != "" {
		if configFile, err := GetConfigFile(env); err == nil {
			return loader.LoadFromFile(configFile)
		}
	}

	return loader.Load()
}

func SaveConfig(config *Config, path string) error {

	return fmt.Errorf("not implemented yet")
}
