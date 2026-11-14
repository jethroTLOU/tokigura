package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/jethroTLOU/tokigura/internal/server/middleware"
	"go.uber.org/fx"
)

type Config struct {
	Port         string
	AllowOrigins []string
	AppName      string
}

func NewServer(lc fx.Lifecycle, cfg *Config) *http.ServeMux {
	mux := http.NewServeMux()

	slog.Info("Registering cors", "allowed origins", strings.Join(cfg.AllowOrigins, ", "))
	handler := middleware.CORS(cfg.AllowOrigins)(mux)
	handler = middleware.Recover(handler)

	srv := &http.Server{
		Addr:              cfg.Port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			slog.Info("Starting HTTP server...", "addr", cfg.Port)

			go func() {
				if err := srv.ListenAndServe(); err != nil &&
					!errors.Is(err, http.ErrServerClosed) {
					slog.Error("HTTP server error", "error", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			slog.Info("Stopping HTTP server...")

			return srv.Shutdown(ctx)
		},
	})

	return mux
}
