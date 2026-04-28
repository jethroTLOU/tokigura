package jwt

import "go.uber.org/fx"

//nolint:gochecknoglobals // uber fx module
var Module = fx.Module("jwt",
	fx.Provide(NewService),
)
