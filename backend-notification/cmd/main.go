package main

import (
	"log"

	"github.com/Wave-95/boards/backend-notification/clients/boards"
	"github.com/Wave-95/boards/backend-notification/clients/email"
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

	// Initialize clients
	amqp, err := amqp.New(cfg.Amqp.User, cfg.Amqp.Password, cfg.Amqp.Host, cfg.Amqp.Port)
	if err != nil {
		log.Fatalf("Error setting up amqp: %v", err)
	}
	amqp.Declare(queues.Notification, 10000, true) //10s ttl with dlx

	emailClient := email.NewClient(cfg.Email.From, cfg.Email.Password, cfg.Email.Host, cfg.Email.Port)
	boardsClient := boards.NewClient(cfg.BoardsBaseURL)

	taskHandler := handlers.New(emailClient, boardsClient, amqp)
	err = taskHandler.Run()

	if err != nil {
		log.Fatalf("Error running task handler: %v", err)
	}
}
