package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/jethroTLOU/tokigura/internal/core/config"
	"github.com/jethroTLOU/tokigura/internal/core/database"
	"github.com/jethroTLOU/tokigura/internal/core/jwt"
	"github.com/jethroTLOU/tokigura/internal/core/logger"
	"github.com/jethroTLOU/tokigura/internal/server"
	"go.uber.org/fx"
	_ "modernc.org/sqlite"
)

func main() {
	app := fx.New(
		logger.Module,
		config.Module,
		server.Module,
		database.Module,
		jwt.Module,
		fx.Provide(newValidator),
	)

	app.Run()
}

func newValidator() *validator.Validate {
	return validator.New(validator.WithRequiredStructEnabled())
}
