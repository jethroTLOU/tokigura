package database

import (
	"context"
	"database/sql"
	"embed"
	"log/slog"
	"time"

	"github.com/pressly/goose/v3"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"go.uber.org/fx"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type Config struct {
	DSN string `json:"dsn"`
}

func New(lc fx.Lifecycle, cfg Config) (*bun.DB, error) {
	sqldb, err := sql.Open("sqlite", cfg.DSN)
	if err != nil {
		slog.Error("Failed to open SQLite", "error", err)
		return nil, err
	}

	sqldb.SetMaxOpenConns(50)
	sqldb.SetMaxIdleConns(10)
	sqldb.SetConnMaxLifetime(time.Hour)

	if err := sqldb.Ping(); err != nil {
		slog.Error("Failed to ping SQLite", "error", err)
		return nil, err
	}

	db := bun.NewDB(sqldb, sqlitedialect.New())

	slog.Info("Connected to SQLite successfully", "dsn", cfg.DSN)

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			slog.Info("Closing SQLite connection...")
			return db.Close()
		},
	})

	return db, nil
}

func RunMigrations(db *bun.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite"); err != nil {
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

	slog.Debug(
		"Migration",
		"current_db_version", currentVer,
		"files_found_in_embed", len(migrations),
	)

	for _, m := range migrations {
		slog.Debug("Available migration", "version", m.Version, "source", m.Source)
	}

	if err = goose.Up(db.DB, "migrations"); err != nil {
		return err
	}

	return nil
}
