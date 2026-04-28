package config

import "go.uber.org/fx"

//nolint:gochecknoglobals // uber fx module
var Module = fx.Module("config",
	fx.Provide(New),
)
