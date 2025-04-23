package service

import (
	"context"
	"manager/internal/dto"
	"manager/internal/model"
	"manager/internal/rabbitmq"

	"go.uber.org/zap"
)

type SubTaskQueueServiceInterface interface {
	SendSubTask(context.Context, *model.SubTask) error
}

type SubTaskQueueService struct {
	logger          *zap.Logger
	rabbitmqManager rabbitmq.RabbitMQManager
}

func NewSubTaskQueueService(logger *zap.Logger, rabbitmqManager rabbitmq.RabbitMQManager) SubTaskQueueServiceInterface {
	return &SubTaskQueueService{
		logger:          logger,
		rabbitmqManager: rabbitmqManager,
	}
}

func (s *SubTaskQueueService) SendSubTask(ctx context.Context, subTask *model.SubTask) error {
	subTaskDto := dto.SubTask{
		TaskID:     subTask.TaskID,
		Hash:       subTask.Hash,
		Length:     subTask.Length,
		Alphabet:   subTask.Alphabet,
		PartNumber: subTask.PartNumber,
		PartCount:  subTask.PartCount,
		ID:         subTask.ID,
	}

	return s.rabbitmqManager.PublishSubtask(ctx, &subTaskDto)
}
