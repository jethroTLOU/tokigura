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

type authCfg struct {
	JWTSecret           string `json:"jwtSecret"`
	SecureCookies       bool   `json:"secureCookies"`
	DeviceSalt          string `json:"deviceSalt"`
	AccessTokeDuration  int    `json:"accessTokeDuration"`
	RefreshTokeDuration int    `json:"refreshTokeDuration"`
}

type Config struct {
	Auth       *authCfg `json:"auth"`
	Database   string   `json:"database"`
	Port       int      `json:"port"`
	Debug      bool     `json:"debug"`
	UploadsDir string   `json:"uploadsDir"`
}

const (
	defaultPort                 = 6969
	defaultAccessTokenDuration  = 15
	defaultRefreshTokenDuration = 30
	randomStringBytes           = 32
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
	if IsDevMode() {
		return "tokigura_dev"
	}
	return "tokigura"
}

func IsDevMode() bool {
	return os.Getenv("GO_DEV") != ""
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
		Database:   fmt.Sprintf("file:%s?cache=shared&mode=rwc&_fk=1", filepath.ToSlash(dbPath)),
		Port:       defaultPort,
		Debug:      false,
		UploadsDir: uploadsDir,
		Auth: &authCfg{
			JWTSecret:           generateRandomString(),
			SecureCookies:       false,
			DeviceSalt:          generateRandomString(),
			AccessTokeDuration:  defaultAccessTokenDuration,
			RefreshTokeDuration: defaultRefreshTokenDuration,
		},
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
