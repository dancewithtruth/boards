package db

import (
	"testing"

	"github.com/Wave-95/boards/server/internal/config"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	// load .env
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}
	// load env vars into config
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Error loading config: %v", err)
	}

	db, err := Connect(cfg.DatabaseConfig)
	defer db.Close()

	assert.NotNil(t, db)
	assert.NoError(t, err)
}
