package consumer

import (
	"context"
	"encoding/json"
	"worker/internal/config"
	"worker/internal/dto"
	"worker/internal/service"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type ProgressConsumer struct {
	logger  *zap.Logger
	conn    *amqp.Connection
	config  *config.Config
	service *service.BruteForceService
}

func NewProgressConsumer(logger *zap.Logger, conn *amqp.Connection, cfg *config.Config, svc *service.BruteForceService) *ProgressConsumer {
	return &ProgressConsumer{
		logger:  logger,
		conn:    conn,
		config:  cfg,
		service: svc,
	}
}

func (c *ProgressConsumer) Consume(ctx context.Context) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	queueName := c.config.ProgressQueueName

	msgs, err := ch.Consume(
		queueName,
		c.config.WorkerName+"-progress",
		false, false, false, false,
		nil,
	)
	if err != nil {
		return err
	}

	c.logger.Info("[📡] ProgressConsumer started", zap.String("queue", queueName))

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("ProgressConsumer shutting down")
			return nil
		case msg := <-msgs:
			var req dto.ProgressRequest
			if err := json.Unmarshal(msg.Body, &req); err != nil {
				c.logger.Error("Invalid progress request", zap.Error(err))
				msg.Nack(false, false)
				continue
			}

			progress := c.service.GetProgress(req.TaskID)

			response := dto.ProgressResponse{
				TaskID:   req.TaskID,
				Progress: progress,
			}

			body, err := json.Marshal(response)
			if err != nil {
				c.logger.Error("Failed to encode response", zap.Error(err))
				msg.Nack(false, false)
				continue
			}

			if msg.ReplyTo != "" {
				err = ch.PublishWithContext(ctx,
					"",          // default exchange
					msg.ReplyTo, // куда отвечать
					false, false,
					amqp.Publishing{
						ContentType:   "application/json",
						CorrelationId: msg.CorrelationId,
						Body:          body,
					})
				if err != nil {
					c.logger.Error("Failed to publish progress response", zap.Error(err))
					msg.Nack(false, false)
					continue
				}
			} else {
				c.logger.Warn("ReplyTo is empty, skipping response")
			}

			msg.Ack(false)
		}
	}
}
