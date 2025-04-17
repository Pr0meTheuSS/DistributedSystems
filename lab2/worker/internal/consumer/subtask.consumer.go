package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"worker/internal/config"
	"worker/internal/dto"
	"worker/internal/service"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type SubTaskConsumer struct {
	service *service.BruteForceService
	conn    *amqp.Connection
	config  *config.Config
	logger  *zap.Logger
}

func NewSubTaskConsumer(logger *zap.Logger, conn *amqp.Connection, cfg *config.Config, service *service.BruteForceService) *SubTaskConsumer {
	return &SubTaskConsumer{
		conn:    conn,
		config:  cfg,
		service: service,
		logger:  logger,
	}
}

func (c *SubTaskConsumer) Consume(ctx context.Context) error {
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
		c.config.QueueName,
		c.config.WorkerName,
		false, // autoAck=false → вручную подтверждаем
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	c.logger.Info("[🚀] Worker started:", zap.String("worker-unique-name", c.config.WorkerName))

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Shutting down consumer")
			return nil
		case msg := <-msgs:
			var task dto.SubTask
			if err := json.Unmarshal(msg.Body, &task); err != nil {
				c.logger.Error("[❌] Failed to decode message:", zap.Error(err))
				msg.Nack(false, false)
				continue
			}

			c.logger.Info("[🚀] Received task:", zap.Any("task", task))

			go c.service.Crack(ctx, &task)
			for i := 0; i < 50; i++ {
				progress := c.service.GetProgress(task.TaskID)
				time.Sleep(time.Second)
				fmt.Println("Progress: ", progress)
			}
			msg.Ack(false)
		}
	}
}
