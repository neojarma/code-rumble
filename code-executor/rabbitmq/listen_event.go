package event_rabbitmq

import (
	"context"
	"encoding/json"
	"executor/connection"
	"executor/entity"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ListenEvent(ctx context.Context, queueName string, function func(*entity.Submission)) error {
	conn, err := connection.GetConnectionRabbitMq()
	if err != nil {
		return err
	}

	rabbitChan, err := conn.Channel()
	if err != nil {
		return err
	}

	q, err := declareQueue(rabbitChan, queueName)
	if err != nil {
		return err
	}

	msgs, err := consumeEvent(rabbitChan, q.Name)
	if err != nil {
		return err
	}

	err = doTask(msgs, function)
	if err != nil {
		return err
	}

	return nil
}

func doTask(msgs <-chan amqp.Delivery, function func(*entity.Submission)) error {
	log.Println("ready to listen event")
	for v := range msgs {
		payload := new(entity.Submission)
		err := json.Unmarshal(v.Body, payload)
		if err != nil {
			return err
		}

		log.Println("do event")
		function(payload)
	}

	return nil
}

func declareQueue(rabbitChan *amqp.Channel, queueName string) (amqp.Queue, error) {
	return rabbitChan.QueueDeclare(
		queueName,
		true, //durable
		false,
		false,
		false,
		nil,
	)
}

func consumeEvent(rabbitChan *amqp.Channel, queueName string) (<-chan amqp.Delivery, error) {
	return rabbitChan.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}
