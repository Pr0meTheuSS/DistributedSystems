package di

import (
	"errors"
	"log"
	"worker/internal/config"
	"worker/internal/consumer"
	"worker/internal/producer"
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

	answersProducer, err := producer.NewProducer(conn, logger, "answers_queue")
	if err != nil {
		logger.Error("Failed to create producer", zap.Error(err))
		panic(err.Error())
	}

	bruteForceService := service.NewBruteForceService(logger)
	if bruteForceService == nil {
		log.Fatal(errors.New("cannot create brute force service component"))
	}

	subTasksConsumer := consumer.NewSubTaskConsumer(logger, conn, cfg, bruteForceService, *answersProducer)

	return &AppContainer{
		Config:   cfg,
		RabbitMQ: conn,
		Consumer: subTasksConsumer,
	}, nil
}
