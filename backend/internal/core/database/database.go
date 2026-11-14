package database

import (
	"context"
	"database/sql"
	"embed"
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

	if err := sqldb.PingContext(context.Background()); err != nil {
		slog.Error("Failed to ping SQLite", "error", err)
		return nil, err
	}

	bunDB := bun.NewDB(sqldb, sqlitedialect.New())
	slog.Info("Connected to SQLite successfully", "dsn", cfg.DSN)

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			slog.Info("Closing SQLite connection...")
			return bunDB.Close()
		},
	})

	return &DB{DB: bunDB}, nil
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

func RunMigrations(db *DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite"); err != nil {
		return err
	}

	currentVer, err := goose.GetDBVersion(db.DB.DB)
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

	if err := goose.Up(db.DB.DB, "migrations"); err != nil {
		return err
	}

	return nil
}
