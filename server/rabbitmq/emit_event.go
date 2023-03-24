package event_rabbitmq

import (
	"context"

	"github.com/rabbitmq/amqp091-go"
)

func EmitEvent(ctx context.Context, body []byte, queueName string, conn *amqp091.Connection) error {
	rabbitChan, err := conn.Channel()
	if err != nil {
		return err
	}

	q, err := declareQueue(rabbitChan, queueName)
	if err != nil {
		return err
	}

	return publishMessage(ctx, rabbitChan, body, q.Name)
}

func declareQueue(rabbitChan *amqp091.Channel, queueName string) (amqp091.Queue, error) {
	return rabbitChan.QueueDeclare(
		queueName,
		true, //durable
		false,
		false,
		false,
		nil,
	)
}

func publishMessage(ctx context.Context, rabbitChan *amqp091.Channel, body []byte, queueName string) error {
	return rabbitChan.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
}
