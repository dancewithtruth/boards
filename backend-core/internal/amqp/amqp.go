package amqp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Wave-95/boards/backend-core/internal/config"
	rabbitmq "github.com/rabbitmq/amqp091-go"
)

const (
	TaskEmailInvites = "task_emal_invites"
)

type Message struct {
	Task    string `json:"task"`
	Payload any    `json:"payload"`
}

type Amqp interface {
	Publish(task string, v any) error
}

type amqpImp struct {
	conn *rabbitmq.Connection
	ch   *rabbitmq.Channel
	q    *rabbitmq.Queue
}

// New creates an Amqp implemented with RabbitMQ. It connects to a broker, opens a channel, and
// declares a notification queue that is configured with durable messages.
func New(cfg config.AmqpConfig) (*amqpImp, error) {
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
		"notification_queue", // name
		true,                 // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	return &amqpImp{conn: conn, ch: ch, q: &q}, nil
}

// Publish publishes a new durable message to the work queue to be processed
// by a consumer.
func (a *amqpImp) Publish(task string, v any) error {
	msg := Message{Task: task, Payload: v}
	bytes, err := json.Marshal(msg)
	err = a.ch.PublishWithContext(context.Background(),
		"",       // exchange
		a.q.Name, // routing key
		false,    // mandatory
		false,
		rabbitmq.Publishing{
			DeliveryMode: rabbitmq.Persistent,
			ContentType:  "text/plain",
			Body:         bytes,
		})

	return err
}
