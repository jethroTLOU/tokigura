package server

import (
	"context"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/jethroTLOU/tokigura/internal/config"
	"go.uber.org/fx"
)

type Config struct {
	Port string
}

func NewServer(lc fx.Lifecycle, cfg *Config) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: config.GetAppName(),
	})

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			slog.Info("Starting Fiber server...")

			go func() {
				if err := app.Listen(cfg.Port); err != nil {
					slog.Error("Fiber server error", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(_ context.Context) error {
			slog.Info("Stopping Fiber server...")
			return app.Shutdown()
		},
	})

	return app
}

func RegisterRoutes(app *fiber.App) {
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
}
