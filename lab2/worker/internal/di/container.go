package di

import (
	"errors"
	"log"
	"worker/internal/config"
	"worker/internal/consumer"
	"worker/internal/rabbitmq"
	"worker/internal/service"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type AppContainer struct {
	Config   *config.Config
	RabbitMQ *amqp.Connection
	Consumer *consumer.SubTaskConsumer
}

func InitContainer() (*AppContainer, error) {
	cfg := config.LoadConfig()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err.Error())
	}

	conn, err := rabbitmq.NewConnection(cfg.RabbitMQURL)
	if err != nil {
		return nil, err
	}

	bruteForceService := service.NewBruteForceService(logger)
	if bruteForceService == nil {
		log.Fatal(errors.New("cannot create brute force service component"))
	}

	consumer := consumer.NewSubTaskConsumer(logger, conn, cfg, bruteForceService)

	return &AppContainer{
		Config:   cfg,
		RabbitMQ: conn,
		Consumer: consumer,
	}, nil
}
