package auth

import (
	"github.com/jethroTLOU/tokigura/internal/identity/session"
	"github.com/jethroTLOU/tokigura/internal/identity/user"
	"github.com/jethroTLOU/tokigura/internal/server"
	"go.uber.org/fx"
)

var Module = fx.Module("auth", fx.Provide(
	func(s *session.Service, u *user.Service) (sessions, userService) {
		return s, u
	},
	fx.Annotate(
		NewService,
		fx.As(new(service)),
	),
	fx.Annotate(
		NewHandler,
		fx.As(new(server.Route)),
		fx.ResultTags(`group:"routes"`),
	),
))
