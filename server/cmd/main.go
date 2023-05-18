package main

import (
	"log"

	"github.com/Wave-95/boards/server/db"
	"github.com/Wave-95/boards/server/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	// load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// load env vars into config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	// run migrations
	err = db.Migrate(cfg.DatabaseConfig)
	if err != nil {
		log.Fatalf("Error running db migrations: %v", err)
	}
	// connect to db
	db, err := db.Connect(cfg.DatabaseConfig)
	if err != nil {
		log.Fatalf("Error connecting to db: %v", err)
	}
	defer db.Close()
}
