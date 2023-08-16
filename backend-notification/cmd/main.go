package main

import (
	"log"

	"github.com/Wave-95/boards/backend-notification/internal/amqp"
	"github.com/Wave-95/boards/backend-notification/internal/config"
)

func main() {
	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	amqp, err := amqp.New(cfg.Amqp)
	if err != nil {
		log.Fatalf("Error setting up amqp: %v", err)
	}

	amqp.Consume()
}
