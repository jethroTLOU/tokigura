package database

import (
	"context"

	"go.uber.org/fx"
)

var Module = fx.Module("database",
	fx.Provide(New),
	fx.Invoke(func(lc fx.Lifecycle, db *DB) {
		lc.Append(fx.Hook{
			OnStart: func(_ context.Context) error {
				return RunMigrations(db)
			},
		})
	}),
)
