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
	JWTSecret            string `json:"jwtSecret"`
	SecureCookies        bool   `json:"secureCookies"`
	AccessTokenDuration  int    `json:"accessTokenDuration"`
	RefreshTokenDuration int    `json:"refreshTokenDuration"`
}

type serverCfg struct {
	Port         int      `json:"port"`
	AllowOrigins []string `json:"allowOrigins"`
}

type Config struct {
	Auth       authCfg   `json:"auth"`
	Server     serverCfg `json:"server"`
	Database   string    `json:"database"`
	Debug      bool      `json:"debug"`
	UploadsDir string    `json:"uploadsDir"`
}

var defaultCors = []string{"http://localhost:5173", "http://127.0.0.1:5173"}

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
	uploadsDir, _ := xdg.DataFile(filepath.Join(appName, "uploads"))

	if err := os.MkdirAll(filepath.Dir(configPath), 0o750); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(uploadsDir, 0o750); err != nil {
		return nil, err
	}

	fileData, err := os.ReadFile(configPath) //nolint:gosec // safe
	if err != nil {
		if os.IsNotExist(err) {
			return createDefaultConfig(configPath, dbPath, uploadsDir)
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(fileData, &cfg); err != nil {
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
	_, _ = rand.Read(b)

	return base64.StdEncoding.EncodeToString(b)
}

func createDefaultConfig(configPath, dbPath, uploadsDir string) (*Config, error) {
	dsn := fmt.Sprintf(
		"file:%s?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=synchronous(NORMAL)",
		filepath.ToSlash(dbPath),
	)

	cfg := &Config{
		Database:   dsn,
		Debug:      false,
		UploadsDir: uploadsDir,
		Server: serverCfg{
			Port:         defaultPort,
			AllowOrigins: defaultCors,
		},
		Auth: authCfg{
			JWTSecret:            generateRandomString(),
			SecureCookies:        false,
			AccessTokenDuration:  defaultAccessTokenDuration,
			RefreshTokenDuration: defaultRefreshTokenDuration,
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
