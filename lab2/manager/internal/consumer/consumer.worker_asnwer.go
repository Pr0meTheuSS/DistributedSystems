package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"manager/internal/config"
	"manager/internal/dto"
	"manager/internal/rabbitmq"
	"manager/internal/service"
	"time"

	"go.uber.org/zap"
)

type WorkerResponseConsumer struct {
	service service.CrackHashServiceInterface
	rabbit  *rabbitmq.RabbitMQManager
	config  *config.Config
	logger  *zap.Logger
}

func NewWorkerResponseConsumer(logger *zap.Logger, rabbit *rabbitmq.RabbitMQManager, cfg *config.Config, service service.CrackHashServiceInterface) *WorkerResponseConsumer {
	return &WorkerResponseConsumer{
		rabbit:  rabbit,
		config:  cfg,
		service: service,
		logger:  logger,
	}
}

func (c *WorkerResponseConsumer) Consume(ctx context.Context) error {
	for {
		if err := c.consumeOnce(ctx); err != nil {
			c.logger.Error("Consumer failed, will retry", zap.Error(err))
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(5 * time.Second):
			}
		} else {
			break
		}
	}
	return nil
}

func (c *WorkerResponseConsumer) consumeOnce(ctx context.Context) error {
	ch, err := c.rabbit.GetChannel()
	if err != nil {
		return fmt.Errorf("cannot open channel: %w", err)
	}
	defer ch.Close()

	err = ch.Qos(1, 0, false)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	msgs, err := ch.Consume(
		c.config.GetWorkerResponseQueue(),
		"manager_consumer",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	c.logger.Info("WorkerResponseConsumer started")

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Shutting down WorkerResponseConsumer")
			return nil
		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("channel closed, reconnect needed")
			}

			var response dto.WorkerResponseDto
			if err := json.Unmarshal(msg.Body, &response); err != nil {
				c.logger.Error("Failed to decode worker response", zap.Error(err))
				msg.Nack(false, false)
				continue
			}

			c.logger.Info("Received worker response", zap.Any("task", response))

			if err := c.service.UpdateProgressOrResult(ctx, response); err != nil {
				c.logger.Error("Failed to handle worker response", zap.Error(err))
				msg.Nack(false, false)
				continue
			}

			msg.Ack(false)
		}
	}
}
