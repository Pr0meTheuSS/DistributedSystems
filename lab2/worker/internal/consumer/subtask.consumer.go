package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"worker/internal/config"
	"worker/internal/dto"
	"worker/internal/model"
	"worker/internal/producer"
	"worker/internal/service"

	"slices"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type SubTaskConsumer struct {
	service  *service.BruteForceService
	conn     *amqp.Connection
	producer *producer.Producer

	statesToSendLetter []model.TaskState
	config             *config.Config
	logger             *zap.Logger
}

func NewSubTaskConsumer(logger *zap.Logger, conn *amqp.Connection, cfg *config.Config, service *service.BruteForceService, producer *producer.Producer) *SubTaskConsumer {
	return &SubTaskConsumer{
		conn:     conn,
		config:   cfg,
		service:  service,
		logger:   logger,
		producer: producer,
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
		false,
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
			// TODO: tryResendToQueue()
			if len(c.statesToSendLetter) != 0 {
				fmt.Println(c.statesToSendLetter)
				for i := 0; i < len(c.statesToSendLetter); i++ {
					state := c.statesToSendLetter[i]

					if err := c.producer.SendToQueue(state); err == nil {
						c.statesToSendLetter = slices.Delete(c.statesToSendLetter, i, i+1)
						i--
					}
				}
			}

			var task dto.SubTask
			if err := json.Unmarshal(msg.Body, &task); err != nil {
				msg.Nack(false, false)
				continue
			}

			c.logger.Info("[🚀] Received task:", zap.Any("task", task))

			go c.service.Crack(ctx, &task)

			state := c.service.GetTaskState(task.ID)
			state.TaskID = task.TaskID

			for state.Status != "READY" {
				state = c.service.GetTaskState(task.ID)
				state.TaskID = task.TaskID

				c.producer.SendToQueue(state)

				time.Sleep(time.Second)
				fmt.Println(c.statesToSendLetter)
			}

			if err := c.producer.SendToQueue(state); err != nil {
				c.logger.Error("Error while send sub task state to queue", zap.Error(err))
				c.statesToSendLetter = append(c.statesToSendLetter, state)
			}

			fmt.Println(c.statesToSendLetter)
			msg.Ack(false)
		}
	}
}
