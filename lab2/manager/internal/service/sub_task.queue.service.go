package service

import (
	"context"
	"manager/internal/dto"
	"manager/internal/rabbitmq"

	"go.uber.org/zap"
)

type SubTaskQueueServiceInterface interface {
	SendSubTask(context.Context, *dto.SubTask) error
}

type SubTaskQueueService struct {
	logger *zap.Logger
	rabbitmqManager rabbitmq.RabbitMQManager
}

func NewSubTaskQueueService(logger *zap.Logger, rabbitmqManager rabbitmq.RabbitMQManager) SubTaskQueueServiceInterface {
	return &SubTaskQueueService{
		logger:          logger,
		rabbitmqManager: rabbitmqManager,
	}
}

func (s *SubTaskQueueService) SendSubTask(ctx context.Context, subTask *dto.SubTask) error {
	return s.rabbitmqManager.PublishSubtask(ctx, subTask)
}
