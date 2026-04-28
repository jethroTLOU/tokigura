package logger

import (
	"log/slog"

	"go.uber.org/fx"
)

//nolint:gochecknoglobals // uber fx module
var Module = fx.Module("logger",
	fx.Provide(New),
	fx.Invoke(func(_ *slog.Logger) {}),
)
