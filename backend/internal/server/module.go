package server

import (
	"fmt"

	"github.com/jethroTLOU/tokigura/internal/config"
	"go.uber.org/fx"
)

//nolint:gochecknoglobals // uber fx module
var Module = fx.Module("server",
	fx.Provide(
		func(cfg *config.Config) *Config {
			return &Config{Port: fmt.Sprintf(":%d", cfg.Port)}
		},
		NewServer,
	),
	fx.Invoke(RegisterRoutes),
)
