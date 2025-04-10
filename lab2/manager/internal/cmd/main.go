package main

import (
	"log"
	"manager/internal/app"
	"manager/internal/config"
	"manager/internal/router"

	"go.uber.org/zap"
)

func main() {
	config := config.NewConfig()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Cannot init zap Logger by call NewProduction() method, err: %s", err.Error())
	}
	defer logger.Sync()

	router := router.NewRouter(logger)

	app := app.NewApp(config, logger, router)

	log.Fatal(app.Run().Error())
}
