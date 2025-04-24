package main

import (
	"log"
	"manager/internal/app"
	"manager/internal/di"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	RabbitMQURL := getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")

	conn, err := amqp.Dial(RabbitMQURL)
	if err != nil {
		logger.Fatal("Failed to connect to RabbitMQ", zap.Error(err))
	}
	defer conn.Close()

	container, err := di.NewContainer(logger, RabbitMQURL)
	if err != nil {
		logger.Fatal("Failed to initialize container", zap.Error(err))
	}

	app := app.NewApp(container)
	if err := app.Run(); err != nil {
		logger.Fatal("Application failed", zap.Error(err))
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Printf("[⚠️] %s not set, using default: %s", key, fallback)
		return fallback
	}
	return val
}
