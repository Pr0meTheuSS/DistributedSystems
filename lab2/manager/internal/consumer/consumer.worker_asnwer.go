package consumer

import (
	"context"
	"encoding/json"
	"manager/internal/config"
	"manager/internal/dto"
	"manager/internal/service"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type WorkerResponseConsumer struct {
	service service.CrackHashServiceInterface
	conn    *amqp.Connection
	config  *config.Config
	logger  *zap.Logger
}

func NewWorkerResponseConsumer(logger *zap.Logger, conn *amqp.Connection, cfg *config.Config, service service.CrackHashServiceInterface) *WorkerResponseConsumer {
	return &WorkerResponseConsumer{
		conn:    conn,
		config:  cfg,
		service: service,
		logger:  logger,
	}
}

func (c *WorkerResponseConsumer) Consume(ctx context.Context) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.Qos(1, 0, false)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		c.config.WorkerResponseQueue,
		"manager_consumer",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	c.logger.Info("[📥] WorkerResponseConsumer started")

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Shutting down WorkerResponseConsumer")
			return nil
		case msg := <-msgs:
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
