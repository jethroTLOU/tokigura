package database

import (
	"context"
	"embed"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"go.uber.org/fx"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type Config struct {
	DSN string `json:"dsn"`
}

func New(lc fx.Lifecycle, cfg *Config, logger *slog.Logger) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite", cfg.DSN)
	if err != nil {
		logger.Error("Failed to connect to SQLite", "error", err)
		return nil, err
	}

	logger.Info("Connected to SQLite successfully", "dsn", cfg.DSN)

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			slog.Info("Closing SQLite connection...")
			return db.Close()
		},
	})

	return db, nil
}

func RunMigrations(db *sqlx.DB, logger *slog.Logger) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	currentVer, err := goose.GetDBVersion(db.DB)
	if err != nil {
		return err
	}

	migrations, err := goose.CollectMigrations("migrations", 0, goose.MaxVersion)
	if err != nil {
		return err
	}

	logger.Debug("Migration", "current_db_version", currentVer, "files_found_in_embed", len(migrations))

	for _, m := range migrations {
		logger.Debug("Available migration", "version", m.Version, "source", m.Source)
	}

	if err = goose.Up(db.DB, "migrations"); err != nil {
		return err
	}

	return nil
}
