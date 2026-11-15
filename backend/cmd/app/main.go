package main

import (
	"fmt"

	"github.com/jethroTLOU/tokigura/internal/core/config"
	"github.com/jethroTLOU/tokigura/internal/core/database"
	"github.com/jethroTLOU/tokigura/internal/core/jwt"
	"github.com/jethroTLOU/tokigura/internal/core/logger"
	"github.com/jethroTLOU/tokigura/internal/identity/auth"
	"github.com/jethroTLOU/tokigura/internal/identity/session"
	"github.com/jethroTLOU/tokigura/internal/identity/user"
	"github.com/jethroTLOU/tokigura/internal/server"
	"go.uber.org/fx"
)

type moduleConfigs struct {
	fx.Out

	Server   server.Config
	Database database.Config
	Jwt      jwt.Config
}

func main() {
	app := fx.New(
		logger.Module,
		config.Module,
		fx.Provide(provideModuleConfigs),
		server.Module,
		database.Module,
		jwt.Module,
		auth.Module,
		session.Module,
		user.Module,
	)

	app.Run()
}

func provideModuleConfigs(cfg *config.Config) moduleConfigs {
	return moduleConfigs{
		Server: server.Config{
			Port:         fmt.Sprintf(":%d", cfg.Server.Port),
			AppName:      config.GetAppName(),
			AllowOrigins: cfg.Server.AllowOrigins,
		},
		Jwt: jwt.Config{
			JWTSecret:            cfg.Auth.JWTSecret,
			AccessTokenDuration:  cfg.Auth.AccessTokenDuration,
			RefreshTokenDuration: cfg.Auth.RefreshTokenDuration,
		},
		Database: database.Config{
			DSN: cfg.Database,
		},
	}
}
