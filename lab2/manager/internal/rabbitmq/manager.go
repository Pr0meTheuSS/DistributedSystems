package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"manager/internal/dto"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type RabbitMQManager struct {
	conn          *amqp.Connection
	logger        *zap.Logger
	subTasksQueue amqp.Queue
	answersQueue  amqp.Queue
}

func NewRabbitMQManager(conn *amqp.Connection, logger *zap.Logger) (*RabbitMQManager, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("cannot open channel: %w", err)
	}
	defer ch.Close()

	err = ch.Confirm(false)
	if err != nil {
		return nil, fmt.Errorf("cannot add confirmation: %w", err)
	}

	subtaskQueue, err := ch.QueueDeclare(
		"subtask_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create queue: %w", err)
	}

	answersQueue, err := ch.QueueDeclare(
		"answers_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create queue: %w", err)
	}

	return &RabbitMQManager{
		conn:          conn,
		logger:        logger,
		subTasksQueue: subtaskQueue,
		answersQueue:  answersQueue,
	}, nil
}

func (m *RabbitMQManager) PublishSubtask(ctx context.Context, subTask *dto.SubTask) error {
	ch, err := m.conn.Channel()
	if err != nil {
		m.logger.Error("Failed to open channel", zap.Error(err))
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer func() {
		if cerr := ch.Close(); cerr != nil {
			m.logger.Warn("Failed to close channel", zap.Error(cerr))
		}
	}()

	err = ch.Confirm(false)
	if err != nil {
		m.logger.Error("Failed to enable publisher confirms", zap.Error(err))
		return fmt.Errorf("failed to enable confirms: %w", err)
	}

	confirmations := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	body, err := json.Marshal(subTask)
	if err != nil {
		m.logger.Error("Failed to serialize task to JSON", zap.Error(err))
		return fmt.Errorf("failed to marshal subtask: %w", err)
	}

	err = ch.PublishWithContext(ctx,
		"",
		m.subTasksQueue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		m.logger.Error("Failed to publish message", zap.Error(err))
		return fmt.Errorf("failed to publish: %w", err)
	}

	select {
	case confirm := <-confirmations:
		if confirm.Ack {
			m.logger.Info("Message confirmed by RabbitMQ")
			return nil
		}

		m.logger.Warn("Message was not acknowledged by RabbitMQ (nack)")
		return fmt.Errorf("message not acknowledged")
	case <-time.After(5 * time.Second):

		m.logger.Error("Timeout waiting for RabbitMQ confirmation")
		return fmt.Errorf("confirmation timeout")
	}
}
