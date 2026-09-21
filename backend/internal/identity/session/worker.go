package session

import (
	"context"
	"log/slog"
	"time"

	"go.uber.org/fx"
)

func startCleanupWorker(lc fx.Lifecycle, s *Service, logger *slog.Logger) {
	ctx, cancel := context.WithCancel(context.Background())

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			logger.Info("Starting session cleanup worker...")

			go func() {
				ticker := time.NewTicker(24 * time.Hour)
				defer ticker.Stop()

				for {
					select {
					case <-ticker.C:
						if err := s.DeleteExpired(ctx); err != nil {
							logger.Error("Failed to clean up expired sessions", "error", err)
						} else {
							logger.Debug("Successfully cleaned up expired sessions")
						}
					case <-ctx.Done():
						return
					}
				}
			}()
			return nil
		},
		OnStop: func(_ context.Context) error {
			logger.Info("Stopping session cleanup worker...")
			cancel()
			return nil
		},
	})
}
