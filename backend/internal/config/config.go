package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Database string `json:"database"`
	Port     int    `json:"port"`
	Debug    bool   `json:"debug"`
}

const defaultPort = 6969

func New() (*Config, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("could not get config dir: %w", err)
	}

	appDir := filepath.Join(configDir, GetAppName())
	configPath := filepath.Join(appDir, "config.json")

	cfg := getDefaultConfig(appDir)

	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		if err = os.MkdirAll(appDir, 0o750); err != nil {
			return nil, fmt.Errorf("failed to create config directory: %w", err)
		}
		if err = saveConfig(configPath, cfg); err != nil {
			return nil, err
		}
		return cfg, nil
	}

	fileData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	if err = json.Unmarshal(fileData, cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	updatedData, _ := json.MarshalIndent(cfg, "", "  ")

	if !bytes.Equal(bytes.TrimSpace(fileData), bytes.TrimSpace(updatedData)) {
		if err = saveConfig(configPath, cfg); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func getDefaultConfig(appDir string) *Config {
	dbPath := filepath.Join(appDir, "database.db")
	dbURI := fmt.Sprintf("file:%s?cache=shared&mode=rwc", filepath.ToSlash(dbPath))

	return &Config{
		Database: dbURI,
		Port:     defaultPort,
		Debug:    false,
	}
}

func saveConfig(path string, cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err = os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func GetAppName() string {
	if os.Getenv("GO_DEV") != "" {
		return "tokigura_dev"
	}
	return "tokigura"
}
