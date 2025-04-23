package producer

import (
	"context"
	"encoding/json"
	"time"
	"worker/internal/model"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Producer struct {
	channel   *amqp.Channel
	loggger   *zap.Logger
	queueName string
}

func NewProducer(conn *amqp.Connection, logger *zap.Logger, queueName string) (*Producer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Producer{
		channel:   ch,
		queueName: queueName,
		loggger:   logger,
	}, nil
}

func (p *Producer) SendToQueue(state model.TaskState) error {
	body, err := json.Marshal(state)
	if err != nil {
		return err
	}
	p.loggger.Info("Sent answer", zap.Any("state", state))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return p.channel.PublishWithContext(
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
}
