package db

import (
	"testing"

	"github.com/Wave-95/boards/backend-core/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	cfg, err := config.Load("../.env")
	if err != nil {
		assert.FailNow(t, "Failed to load config which is needed to test Connect.", err)
	}
	db, err := Connect(cfg.DatabaseConfig)
	assert.NoError(t, err)
	defer db.Close()
	assert.NotNil(t, db)
}
