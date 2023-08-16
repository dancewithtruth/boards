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
	AddHandler(task string, handler func(payload interface{}) error)
}

type amqpClient struct {
	conn     *rabbitmq.Connection
	ch       *rabbitmq.Channel
	handlers map[string]func(payload interface{}) error
}

// New creates an Amqp implemented with RabbitMQ. It connects to a broker and opens a channel
func New(user, password, host, port string) (*amqpClient, error) {
	handlers := make(map[string]func(payload interface{}) error)
	connString := fmt.Sprintf("amqp://%v:%v@%v:%v/", user, password, host, port)
	conn, err := rabbitmq.Dial(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to amqp server: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &amqpClient{conn: conn, ch: ch, handlers: handlers}, nil
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
	if err != nil {
		return fmt.Errorf("failed to marshal the message: %w", err)
	}
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

	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			fmt.Println("message:", string(d.Body))
			fmt.Println("handlers:", a.handlers)
			var msg tasks.Message
			err := json.Unmarshal(d.Body, &msg)
			if err != nil {
				fmt.Printf("failed to unmarshal message: %v", err)
				continue
			}
			if handler, ok := a.handlers[msg.Task]; ok {
				err = handler(msg.Payload)
				if err != nil {
					fmt.Printf("failed to handle message: %v", err)
					continue
				}
				d.Ack(false)
			} else {
				fmt.Printf("task does not exist: %v", msg.Task)
			}
		}
	}()

	<-forever

	return nil
}

func (a *amqpClient) AddHandler(task string, handler func(interface{}) error) {
	a.handlers[task] = handler
}
