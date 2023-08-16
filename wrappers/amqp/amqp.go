package amqp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Wave-95/boards/shared/tasks"
	rabbitmq "github.com/rabbitmq/amqp091-go"
)

type Amqp interface {
	Publish(queue string, task string, v any) error
	Consume(queue string) error
}

type amqpClient struct {
	conn *rabbitmq.Connection
	ch   *rabbitmq.Channel
}

// New creates an Amqp implemented with RabbitMQ. It connects to a broker and opens a channel
func New(user, password, host, port string) (*amqpClient, error) {
	connString := fmt.Sprintf("amqp://%v:%v@%v:%v/", user, password, host, port)
	conn, err := rabbitmq.Dial(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to amqp server: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &amqpClient{conn: conn, ch: ch}, nil
}

// Publish publishes a new durable message to the work queue to be processed
// by a consumer.
func (a *amqpClient) Publish(queue string, task string, v any) error {
	q, err := a.ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	msg := tasks.Message{Task: task, Payload: v}
	bytes, err := json.Marshal(msg)
	err = a.ch.PublishWithContext(context.Background(),
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		rabbitmq.Publishing{
			DeliveryMode: rabbitmq.Persistent,
			ContentType:  "text/plain",
			Body:         bytes,
		})

	return err
}

// Consume is a blocking operation that consumes each new message published to
// a queue.
func (a *amqpClient) Consume(queue string) error {
	_, err := a.ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	err = a.ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return fmt.Errorf("failed to set Qos: %w", err)
	}

	msgs, err := a.ch.Consume(
		queue, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
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
