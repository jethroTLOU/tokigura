package config

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

type Config struct {
	Database   string `json:"database"`
	Port       int    `json:"port"`
	Debug      bool   `json:"debug"`
	JWTSecret  string `json:"jwtSecret"`
	DeviceSalt string `json:"deviceSalt"`
	UploadsDir string `json:"uploadsDir"`
}

const (
	defaultPort       = 6969
	randomStringBytes = 32
)

func New() (*Config, error) {
	appName := GetAppName()

	configPath, _ := xdg.ConfigFile(filepath.Join(appName, "config.json"))
	dbPath, _ := xdg.DataFile(filepath.Join(appName, "database.db"))
	uploadsDir, _ := xdg.DataFile(filepath.Join(appName, "uploads/"))

	if err := os.MkdirAll(filepath.Dir(configPath), 0o750); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(uploadsDir, 0o750); err != nil {
		return nil, err
	}

	fileData, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return createDefaultConfig(configPath, dbPath, uploadsDir)
		}
		return nil, err
	}

	var cfg Config
	if err = json.Unmarshal(fileData, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func GetAppName() string {
	if os.Getenv("GO_DEV") != "" {
		return "tokigura_dev"
	}
	return "tokigura"
}

func generateRandomString() string {
	b := make([]byte, randomStringBytes)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Errorf("failed to generate secure random string: %w", err))
	}
	return base64.StdEncoding.EncodeToString(b)
}

func createDefaultConfig(configPath, dbPath, uploadsDir string) (*Config, error) {
	cfg := &Config{
		Database:   fmt.Sprintf("file:%s?cache=shared&mode=rwc", filepath.ToSlash(dbPath)),
		Port:       defaultPort,
		Debug:      false,
		JWTSecret:  generateRandomString(),
		DeviceSalt: generateRandomString(),
		UploadsDir: uploadsDir,
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal default config: %w", err)
	}

	if err = os.WriteFile(configPath, data, 0o600); err != nil {
		return nil, fmt.Errorf("failed to write default config file: %w", err)
	}

	return cfg, nil
}
