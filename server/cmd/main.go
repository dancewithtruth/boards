package main

import (
	"log"

	"github.com/Wave-95/boards/server/db"
	"github.com/Wave-95/boards/server/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// Load env vars into config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	// Run migrations
	err = db.Migrate(cfg.DatabaseConfig)
	if err != nil {
		log.Fatalf("Error running db migrations: %v", err)
	}
}
