package server

import (
	"context"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/fx"
)

type Config struct {
	Port    string
	AppName string
}

func NewServer(lc fx.Lifecycle, cfg *Config, logger *slog.Logger) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      cfg.AppName,
		ErrorHandler: ErrorHandler(logger),
	})

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			logger.Info("Starting Fiber server...")

			go func() {
				if err := app.Listen(cfg.Port); err != nil {
					logger.Error("Fiber server error", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(_ context.Context) error {
			logger.Info("Stopping Fiber server...")
			return app.Shutdown()
		},
	})

	return app
}
