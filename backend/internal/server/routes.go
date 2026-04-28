package server

import (
	"github.com/gofiber/fiber/v3"
	"go.uber.org/fx"
)

type Route interface {
	Register(router fiber.Router)
}

type RegisterRoutesParams struct {
	fx.In

	App    *fiber.App
	Routes []Route `group:"routes"`
}

func RegisterRoutes(p RegisterRoutesParams) {
	api := p.App.Group("/api/v1")

	for _, r := range p.Routes {
		r.Register(api)
	}
}
