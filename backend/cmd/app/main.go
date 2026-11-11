package main

import (
	"fmt"

	"github.com/jethroTLOU/tokigura/internal/core/config"
	"github.com/jethroTLOU/tokigura/internal/core/database"
	"github.com/jethroTLOU/tokigura/internal/core/jwt"
	"github.com/jethroTLOU/tokigura/internal/core/logger"
	"github.com/jethroTLOU/tokigura/internal/core/validation"
	"github.com/jethroTLOU/tokigura/internal/identity/auth"
	"github.com/jethroTLOU/tokigura/internal/identity/session"
	"github.com/jethroTLOU/tokigura/internal/identity/user"
	"github.com/jethroTLOU/tokigura/internal/server"
	"go.uber.org/fx"
	_ "modernc.org/sqlite"
)

func main() {
	app := fx.New(
		logger.Module,
		config.Module,
		fx.Provide(func(cfg *config.Config) database.Config {
			return database.Config{DSN: cfg.Database}
		}),
		fx.Provide(func(cfg *config.Config) jwt.Config {
			return jwt.Config{
				JWTSecret:            cfg.Auth.JWTSecret,
				AccessTokenDuration:  cfg.Auth.AccessTokenDuration,
				RefreshTokenDuration: cfg.Auth.RefreshTokenDuration,
			}
		}),
		fx.Provide(func(cfg *config.Config) *server.Config {
			return &server.Config{
				Port:         fmt.Sprintf(":%d", cfg.Server.Port),
				AppName:      config.GetAppName(),
				AllowOrigins: cfg.Server.AllowOrigins,
			}
		}),
		server.Module,
		database.Module,
		jwt.Module,
		fx.Provide(validation.New),
		auth.Module,
		session.Module,
		user.Module,
	)

	app.Run()
}
