package database

import (
	"context"
	"embed"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/fx"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type Config struct {
	DSN string `json:"dsn"`
}

func New(lc fx.Lifecycle, cfg *Config, logger *slog.Logger) (*pgxpool.Pool, error) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DSN)
	if err != nil {
		logger.Error("Failed to create database pool.", "error", err)
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		logger.Error("Failed to ping database.", "error", err)
		return nil, err
	}

	logger.Info("Connected to PostgreSQL successfully", "dsn", cfg.DSN)

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			logger.Info("Closing PostgreSQL connection pool...")
			pool.Close()
			return nil
		},
	})

	return pool, nil
}

func RunMigrations(pool *pgxpool.Pool, logger *slog.Logger) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	sqldb := stdlib.OpenDB(*pool.Config().ConnConfig)
	defer sqldb.Close()

	currentVer, err := goose.GetDBVersion(sqldb)
	if err != nil {
		return err
	}

	migrations, err := goose.CollectMigrations("migrations", 0, goose.MaxVersion)
	if err != nil {
		return err
	}

	logger.Debug(
		"Migration",
		"current_db_version",
		currentVer,
		"files_found_in_embed",
		len(migrations),
	)

	for _, m := range migrations {
		logger.Debug("Available migration", "version", m.Version, "source", m.Source)
	}

	if err = goose.Up(sqldb, "migrations"); err != nil {
		return err
	}

	return nil
}
