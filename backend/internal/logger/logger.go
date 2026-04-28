package logger

import (
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
	"github.com/jethroTLOU/tokigura/internal/config"
)

func New(cfg *config.Config) *slog.Logger {
	logLevel := log.InfoLevel
	if cfg.Debug {
		logLevel = log.DebugLevel
	}

	charmHandler := log.NewWithOptions(os.Stderr, log.Options{
		Level:           logLevel,
		ReportTimestamp: true,
	})

	logger := slog.New(charmHandler)

	slog.SetDefault(logger)

	return logger
}
