package amqp

import (
	"fmt"

	"github.com/Wave-95/boards/backend-notification/internal/config"
	"github.com/Wave-95/boards/shared/queues"
	rabbitmq "github.com/rabbitmq/amqp091-go"
)

type Amqp interface {
	Consume() error
}

type amqpClient struct {
	conn *rabbitmq.Connection
	ch   *rabbitmq.Channel
	q    *rabbitmq.Queue
}

// New creates an Amqp implemented with RabbitMQ. It connects to a broker, opens a channel, and
// declares a notification queue that is configured with durable messages.
func New(cfg config.AmqpConfig) (*amqpClient, error) {
	connString := fmt.Sprintf("amqp://%v:%v@%v:%v/", cfg.User, cfg.Password, cfg.Host, cfg.Port)
	conn, err := rabbitmq.Dial(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to amqp server: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	q, err := ch.QueueDeclare(
		queues.Notification, // name
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, fmt.Errorf("failed to set Qos: %w", err)
	}

	return &amqpClient{conn: conn, ch: ch, q: &q}, nil
}

// Publish publishes a new durable message to the work queue to be processed
// by a consumer.
func (a *amqpClient) Consume() error {
	msgs, err := a.ch.Consume(
		a.q.Name, // queue
		"",       // consumer
		false,    // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
	if err != nil {
		return fmt.Errorf("failed to register as consumer: %w", err)
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			// Process each msg
			fmt.Println(d.Body)
			d.Ack(false)
		}
	}()

	<-forever

	return nil
}
