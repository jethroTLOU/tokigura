package database

import (
	"github.com/jethroTLOU/tokigura/internal/core/config"
	"go.uber.org/fx"
)

var Module = fx.Module("database",
	fx.Provide(
		func(cfg *config.Config) *Config {
			return &Config{DSN: cfg.Database}
		},
		New,
	),
	fx.Invoke(RunMigrations),
)
