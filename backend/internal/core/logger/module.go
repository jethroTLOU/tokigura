package logger

import (
	"log/slog"

	"go.uber.org/fx"
)

var Module = fx.Module("logger",
	fx.Provide(New),
	fx.Invoke(func(_ *slog.Logger) {}),
)
