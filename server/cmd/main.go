package main

import (
	"log"
	"net/http"

	"github.com/Wave-95/boards/server/db"
	"github.com/Wave-95/boards/server/internal/api/user"
	"github.com/Wave-95/boards/server/internal/config"
	"github.com/go-chi/chi/v5"
)

func main() {
	// load env vars into config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	// connect to db
	db, err := db.Connect(cfg.DatabaseConfig)
	if err != nil {
		log.Fatalf("Error connecting to db: %v", err)
	}
	defer db.Close()

	r := chi.NewRouter()
	userAPI := user.NewAPI(user.NewService(user.NewRepository(db)))
	userAPI.RegisterHandlers(r)
	http.ListenAndServe(":8080", r)
}
