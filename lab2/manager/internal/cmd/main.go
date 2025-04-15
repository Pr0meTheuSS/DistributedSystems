package main

import (
	"manager/internal/app"
	"manager/internal/di"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		logger.Fatal("Failed to connect to RabbitMQ", zap.Error(err))
	}
	defer conn.Close()

	container, err := di.NewContainer(logger, conn)
	if err != nil {
		logger.Fatal("Failed to initialize container", zap.Error(err))
	}

	app := app.NewApp(container)
	if err := app.Run(); err != nil {
		logger.Fatal("Application failed", zap.Error(err))
	}
}
