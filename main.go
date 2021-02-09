package main

import (
	"github.com/ersmith/mailgun-coding-challenge/config"
	"go.uber.org/zap"
)

// Entrypoint
func main() {
	logger := initLogger()
	defer logger.Sync()
	logger.Infof("Initializing app")

	config := config.Config{
		Logger: logger,
	}
	app := App{}
	app.Initialize(config.DbConfig(), logger)
	defer app.Cleanup()
	app.Run(config.HttpPort())
}

// Initializes the logger
func initLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	return logger.Sugar()
}
