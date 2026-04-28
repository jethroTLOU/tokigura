package main

import (
	"github.com/jethroTLOU/tokigura/internal/config"
	"github.com/jethroTLOU/tokigura/internal/database"
	"github.com/jethroTLOU/tokigura/internal/logger"
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
	)

	app.Run()
}
