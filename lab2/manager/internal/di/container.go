package di

import (
	"log"
	"manager/internal/config"
	"manager/internal/consumer"
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

	CrackHashService       service.CrackHashServiceInterface
	WorkerResponseConsumer *consumer.WorkerResponseConsumer

	PingHandler      *handler.PingHandler
	CrackHashHandler *handler.CrackHashHandler
}

func NewContainer(logger *zap.Logger, rabbitConn *amqp.Connection) (*AppContainer, error) {
	cfg := config.NewConfig()

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

	WorkersAmount := int64(3) // TODO: get from rabbit rest api
	crackHashService := service.NewCrackHashService(logger, subTaskQueueService, requestsRepository, WorkersAmount)

	pingHandler := handler.NewPingHandler(service.NewPingService(logger))
	crackHashHandler := handler.NewCrackHashHandler(crackHashService)

	workerConsumer := consumer.NewWorkerResponseConsumer(
		logger,
		rabbitConn,
		cfg,
		crackHashService,
	)

	return &AppContainer{
		Logger:                 logger,
		Config:                 cfg,
		RabbitManager:          rabbitManager,
		CrackHashService:       crackHashService,
		WorkerResponseConsumer: workerConsumer,
		PingHandler:            &pingHandler,
		CrackHashHandler:       crackHashHandler,
	}, nil
}
