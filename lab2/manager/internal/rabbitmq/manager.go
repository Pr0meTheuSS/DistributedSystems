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
	connStr       string
	conn          *amqp.Connection
	logger        *zap.Logger
	subTasksQueue amqp.Queue
	answersQueue  amqp.Queue
	done          chan bool
	isConnected   bool
}

func (m *RabbitMQManager) GetChannel() (*amqp.Channel, error) {
	if !m.isConnected {
		return nil, fmt.Errorf("RabbitMQ not connected")
	}

	ch, err := m.conn.Channel()
	if err != nil {
		m.logger.Error("Failed to get channel", zap.Error(err))
		return nil, err
	}

	return ch, nil
}

func NewRabbitMQManager(connStr string, logger *zap.Logger) (*RabbitMQManager, error) {
	manager := &RabbitMQManager{
		connStr:     connStr,
		logger:      logger,
		done:        make(chan bool),
		isConnected: false,
	}

	go manager.handleReconnect()
	return manager, nil
}

func (m *RabbitMQManager) handleReconnect() {
	for {
		if m.connect() {
			m.logger.Info("Successfully connected to RabbitMQ")
			break
		}
		m.logger.Warn("Retrying RabbitMQ connection in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
}

func (m *RabbitMQManager) connect() bool {
	conn, err := amqp.Dial(m.connStr)
	if err != nil {
		m.logger.Error("Failed to connect to RabbitMQ", zap.Error(err))
		return false
	}

	ch, err := conn.Channel()
	if err != nil {
		m.logger.Error("Failed to open channel", zap.Error(err))
		return false
	}

	subtaskQueue, err := ch.QueueDeclare("subtask_queue", true, false, false, false, nil)
	if err != nil {
		m.logger.Error("Failed to declare subtask queue", zap.Error(err))
		return false
	}

	answersQueue, err := ch.QueueDeclare("answers_queue", true, false, false, false, nil)
	if err != nil {
		m.logger.Error("Failed to declare answers queue", zap.Error(err))
		return false
	}

	m.conn = conn
	m.subTasksQueue = subtaskQueue
	m.answersQueue = answersQueue
	m.isConnected = true

	go m.monitorConnection()

	return true
}

func (m *RabbitMQManager) monitorConnection() {
	closeChan := m.conn.NotifyClose(make(chan *amqp.Error))
	err := <-closeChan
	m.logger.Warn("RabbitMQ connection closed", zap.Error(err))
	m.isConnected = false
	go m.handleReconnect()
}

func (m *RabbitMQManager) PublishSubtask(ctx context.Context, subTask *dto.SubTask) error {
	ch, err := m.GetChannel()
	if err != nil {
		return fmt.Errorf("cannot get channel: %w", err)
	}
	defer ch.Close()

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
