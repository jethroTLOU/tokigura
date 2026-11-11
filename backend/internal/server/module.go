package server

import (
	"github.com/jethroTLOU/tokigura/internal/server/middleware"
	"go.uber.org/fx"
)

var Module = fx.Module("server",
	fx.Provide(
		NewServer,
		middleware.NewAuthMiddleware,
	),
	fx.Invoke(RegisterRoutes),
)
