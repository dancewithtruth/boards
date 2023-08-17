package main

import (
	"log"

	"github.com/Wave-95/boards/backend-notification/internal/config"
	"github.com/Wave-95/boards/backend-notification/internal/email"
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

	emailClient := email.NewClient("useboards@gmail.com", "wdlrmviaiwalxkkq", "smtp.gmail.com", "587")

	taskHandler := handlers.New(emailClient, amqp)
	err = taskHandler.Run()

	if err != nil {
		log.Fatalf("Error running task handler: %v", err)
	}
}
