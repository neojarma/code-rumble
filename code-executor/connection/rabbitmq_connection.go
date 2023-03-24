package connection

import (
	"errors"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func GetConnectionRabbitMq() (*amqp.Connection, error) {
	// url := os.Getenv("RABBIT_URL")
	url := "amqp://guest:guest@localhost:5672"

	maxTries := 10
	for i := 0; i < maxTries; i++ {
		conn, err := amqp.Dial(url)
		if err == nil {
			log.Println("success connect to rabbitmq")
			return conn, nil
		}

		log.Println("failed to connect to rabbimq, try again in 1 minute")
		time.Sleep(1 * time.Second)
	}

	return nil, errors.New("rabbitmq connection failed after 10 minute")

}
