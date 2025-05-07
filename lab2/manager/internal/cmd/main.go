package main

import (
	"manager/internal/app"
	"manager/internal/di"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	container, err := di.NewContainer(logger)
	if err != nil {
		logger.Fatal("Failed to initialize container", zap.Error(err))
	}

	app := app.NewApp(container)
	if err := app.Run(); err != nil {
		logger.Fatal("Application failed", zap.Error(err))
	}
}
