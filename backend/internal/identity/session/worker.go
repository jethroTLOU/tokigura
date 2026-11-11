package session

import (
	"context"
	"log/slog"
	"time"

	"go.uber.org/fx"
)

func startCleanupWorker(lc fx.Lifecycle, s *Service) {
	ctx, cancel := context.WithCancel(context.Background())

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			slog.Info("Starting session cleanup worker...")

			go func() {
				ticker := time.NewTicker(24 * time.Hour)
				defer ticker.Stop()

				deleteExpired(ctx, s)

				for {
					select {
					case <-ticker.C:
						deleteExpired(ctx, s)
					case <-ctx.Done():
						return
					}
				}
			}()
			return nil
		},
		OnStop: func(_ context.Context) error {
			slog.Info("Stopping session cleanup worker...")
			cancel()
			return nil
		},
	})
}

func deleteExpired(ctx context.Context, s *Service) {
	if err := s.DeleteExpired(ctx); err != nil {
		slog.Error("Failed to clean up expired sessions", "error", err)
	} else {
		slog.Debug("Successfully cleaned up expired sessions")
	}
}
