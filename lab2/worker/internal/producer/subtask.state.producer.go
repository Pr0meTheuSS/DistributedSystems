package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"worker/internal/model"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Producer struct {
	rabbitURI string
	channel   *amqp.Channel
	logger    *zap.Logger
	queueName string
	conn      *amqp.Connection
	mu        sync.RWMutex
}

func NewProducer(conn *amqp.Connection, logger *zap.Logger, rabbitURI, queueName string) (*Producer, error) {
	p := &Producer{
		conn:      conn,
		logger:    logger,
		queueName: queueName,
		rabbitURI: rabbitURI,
	}

	if err := p.setupChannel(); err != nil {
		return nil, err
	}

	go p.handleConnectionLoss()

	return p, nil
}

func (p *Producer) setupChannel() error {
	ch, err := p.conn.Channel()
	if err != nil {
		return err
	}

	_, err = ch.QueueDeclare(
		p.queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	if p.channel != nil {
		p.channel.Close()
	}
	p.channel = ch
	return nil
}

func (p *Producer) handleConnectionLoss() {
	errs := p.conn.NotifyClose(make(chan *amqp.Error))
	for err := range errs {
		p.logger.Error("RabbitMQ connection lost", zap.Error(err))
		time.Sleep(5 * time.Second)

		for {
			newConn, err := amqp.Dial(p.rabbitURI)
			if err != nil {
				p.logger.Error("Error reconnecting to RabbitMQ", zap.Error(err))
				time.Sleep(5 * time.Second)
				continue
			}

			p.conn = newConn
			if err := p.setupChannel(); err == nil {
				p.logger.Info("Reconnected to RabbitMQ successfully")
				break
			} else {
				p.logger.Error("Error setting up channel after reconnecting", zap.Error(err))
				newConn.Close()
				time.Sleep(5 * time.Second)
			}
		}
	}
}

func (p *Producer) SendToQueue(state model.TaskState) error {
	body, err := json.Marshal(state)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p.mu.RLock()
	ch := p.channel
	p.mu.RUnlock()

	if ch == nil || ch.IsClosed() {
		p.logger.Warn("Channel is closed or nil, attempting to recover")
		p.mu.Lock()
		if err := p.setupChannel(); err != nil {
			p.mu.Unlock()
			return fmt.Errorf("failed to reset channel: %w", err)
		}
		ch = p.channel
		p.mu.Unlock()
	}

	err = ch.PublishWithContext(
		ctx,
		"",
		p.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err == nil {
		p.logger.Info("Sent answer", zap.Any("state", state))
	} else {
		p.logger.Error("Send answer error", zap.Error(err))
	}

	return err
}
