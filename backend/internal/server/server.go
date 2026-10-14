package server

import (
	"context"
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"go.uber.org/fx"
)

type Config struct {
	Port         string
	AllowOrigins []string
	AppName      string
}

func NewServer(lc fx.Lifecycle, cfg *Config) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      cfg.AppName,
		ErrorHandler: ErrorHandler(),
	})

	slog.Info("Registering cors", "allowed origins", strings.Join(cfg.AllowOrigins, ", "))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowOrigins,
		AllowCredentials: true,
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods:     []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH", "OPTIONS"},
	}))

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
