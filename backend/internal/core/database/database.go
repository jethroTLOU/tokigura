package database

import (
	"context"
	"database/sql"
	"embed"
	"io/fs"
	"log/slog"

	"github.com/pressly/goose/v3"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"go.uber.org/fx"
	_ "modernc.org/sqlite" // nolint:revive // db
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type txKey struct{}

type Config struct {
	DSN string `json:"dsn"`
}

type DB struct {
	*bun.DB
}

func New(lc fx.Lifecycle, cfg Config) (*DB, error) {
	sqldb, err := sql.Open("sqlite", cfg.DSN)
	if err != nil {
		slog.Error("Failed to open SQLite", "error", err)
		return nil, err
	}

	sqldb.SetMaxOpenConns(1)
	sqldb.SetMaxIdleConns(1)
	sqldb.SetConnMaxLifetime(0)

	bunDB := bun.NewDB(sqldb, sqlitedialect.New())
	db := &DB{DB: bunDB}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := sqldb.PingContext(ctx); err != nil {
				slog.Error("Failed to ping SQLite", "error", err)
				return err
			}
			slog.Info("Connected to SQLite successfully", "dsn", cfg.DSN)

			return runMigrations(ctx, sqldb)
		},
		OnStop: func(_ context.Context) error {
			slog.Info("Closing SQLite connection...")
			return bunDB.Close()
		},
	})

	return db, nil
}

func (d *DB) Conn(ctx context.Context) bun.IDB {
	if tx, ok := ctx.Value(txKey{}).(bun.Tx); ok {
		return tx
	}
	return d.DB
}

func (d *DB) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := ctx.Value(txKey{}).(bun.Tx); ok {
		return fn(ctx)
	}
	return d.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		return fn(context.WithValue(ctx, txKey{}, tx))
	})
}

func runMigrations(ctx context.Context, sqldb *sql.DB) error {
	migrationsFS, err := fs.Sub(embedMigrations, "migrations")
	if err != nil {
		return err
	}

	provider, err := goose.NewProvider(goose.DialectSQLite3, sqldb, migrationsFS)
	if err != nil {
		return err
	}

	results, err := provider.Up(ctx)
	if err != nil {
		return err
	}

	for _, r := range results {
		slog.Info("Applied migration", "version", r.Source.Version, "duration", r.Duration)
	}

	return nil
}
