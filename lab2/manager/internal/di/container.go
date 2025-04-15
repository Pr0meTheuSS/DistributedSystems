package di

import (
	"manager/internal/config"
	"manager/internal/handler"
	"manager/internal/rabbitmq"
	"manager/internal/service"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type AppContainer struct {
	Logger        *zap.Logger
	Config        *config.Config
	RabbitManager *rabbitmq.RabbitMQManager

	CrackHashService service.CrackHashServiceInterface

	PingHandler      *handler.PingHandler
	CrackHashHandler *handler.CrackHashHandler
}

func NewContainer(logger *zap.Logger, rabbitConn *amqp.Connection) (*AppContainer, error) {
	cfg := config.NewConfig() // предположим, ты уже это используешь

	rabbitManager, err := rabbitmq.NewRabbitMQManager(rabbitConn, logger)
	if err != nil {
		return nil, err
	}

	subTaskQueueService := service.NewSubTaskQueueService(logger, *rabbitManager)

	crackHashService := service.NewCrackHashService(logger, subTaskQueueService, 3)
	pingHandler := handler.NewPingHandler(service.NewPingService(logger))
	crackHashHandler := handler.NewCrackHashHandler(crackHashService)

	return &AppContainer{
		Logger:           logger,
		Config:           cfg,
		RabbitManager:    rabbitManager,
		CrackHashService: crackHashService,
		PingHandler:      &pingHandler,
		CrackHashHandler: crackHashHandler,
	}, nil
}
