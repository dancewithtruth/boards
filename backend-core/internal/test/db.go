package test

import (
	"testing"

	"github.com/Wave-95/boards/backend-core/db"
	"github.com/Wave-95/boards/backend-core/internal/config"
)

func DB(t *testing.T) *db.DB {
	// load env vars into config
	cfg, err := config.Load()
	if err != nil {
		t.Errorf("Issue loading config:%v", err)
		t.FailNow()
	}
	// connect to db
	db, err := db.Connect(cfg.DatabaseConfig)
	if err != nil {
		t.Errorf("Issue connecting db:%v", err)
		t.FailNow()
	}
	return db
}
