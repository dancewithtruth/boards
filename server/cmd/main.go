package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// setup server
	r := chi.NewRouter()
	server := http.Server{Addr: ":8080", Handler: buildHandler(r, db)}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %s", err)
		}
	}()

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced shutdown: %s", err)
	}
}

func buildHandler(r chi.Router, db *db.DB) chi.Router {
	// register user handlers
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userAPI := user.NewAPI(userService)
	userAPI.RegisterHandlers(r)

	return r
}
