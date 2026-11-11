package session

import (
	"github.com/jethroTLOU/tokigura/internal/core/jwt"
	"go.uber.org/fx"
)

var Module = fx.Module("session", fx.Provide(
	func(s *jwt.Service) jwtService { return s },
	NewService,
	fx.Annotate(NewRepository, fx.As(new(repository))),
),
	fx.Invoke(startCleanupWorker),
)
