package di

import (
	"log"
	"manager/internal/config"
	"manager/internal/db"
	"manager/internal/handler"
	"manager/internal/rabbitmq"
	"manager/internal/repository"
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

const (
	WorkersAmount = 3
)

func NewContainer(logger *zap.Logger, rabbitConn *amqp.Connection) (*AppContainer, error) {
	cfg := config.NewConfig() // предположим, ты уже это используешь

	rabbitManager, err := rabbitmq.NewRabbitMQManager(rabbitConn, logger)
	if err != nil {
		return nil, err
	}

	subTaskQueueService := service.NewSubTaskQueueService(logger, *rabbitManager)

	db, err := db.NewMongoDatabase("mongodb://root:example@localhost:27017", "bf-service")
	if err != nil {
		log.Fatal(err)
	}
	requestsRepository := repository.NewRequestsMongoRepository(db, logger)

	crackHashService := service.NewCrackHashService(logger, subTaskQueueService, requestsRepository, WorkersAmount)
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
