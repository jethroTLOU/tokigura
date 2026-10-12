package database

import (
	"context"

	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

var Module = fx.Module("database",
	fx.Provide(New),
	fx.Invoke(func(lc fx.Lifecycle, db *bun.DB) {
		lc.Append(fx.Hook{
			OnStart: func(_ context.Context) error {
				return RunMigrations(db)
			},
		})
	}),
)
