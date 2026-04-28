package database

import (
	"github.com/jethroTLOU/tokigura/internal/config"
	"go.uber.org/fx"
)

//nolint:gochecknoglobals // uber fx module
var Module = fx.Module("database",
	fx.Provide(
		func(cfg *config.Config) *Config {
			return &Config{DSN: cfg.Database}
		},
		New,
	),
	fx.Invoke(RunMigrations),
)
