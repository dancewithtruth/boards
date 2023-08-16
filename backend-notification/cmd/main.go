package main

import (
	"log"

	"github.com/Wave-95/boards/backend-notification/constants/queues"
	"github.com/Wave-95/boards/backend-notification/internal/config"
	"github.com/Wave-95/boards/backend-notification/internal/handlers"
	"github.com/Wave-95/boards/wrappers/amqp"
)

func main() {
	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	amqp, err := amqp.New(cfg.Amqp.User, cfg.Amqp.Password, cfg.Amqp.Host, cfg.Amqp.Port)
	if err != nil {
		log.Fatalf("Error setting up amqp: %v", err)
	}

	handlers.Register(amqp)

	err = amqp.Consume(queues.Notification)
	if err != nil {
		log.Fatalf("Error consuming messages from amqp: %v", err)
	}
}
