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
	"github.com/Wave-95/boards/server/internal/api/auth"
	"github.com/Wave-95/boards/server/internal/api/user"
	"github.com/Wave-95/boards/server/internal/config"
	"github.com/Wave-95/boards/server/internal/middleware"
	"github.com/Wave-95/boards/server/pkg/logger"
	"github.com/Wave-95/boards/server/pkg/validator"
	"github.com/go-chi/chi/v5"
)

func main() {
	logger := logger.New()
	validator := validator.New()
	// load env vars into config
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}
	// connect to db
	db, err := db.Connect(cfg.DatabaseConfig)
	if err != nil {
		logger.Fatalf("Error connecting to db: %v", err)
	}
	defer db.Close()

	// setup server
	r := chi.NewRouter()
	server := http.Server{Addr: ":8080", Handler: buildHandler(r, db, logger, validator, cfg)}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Could not start server: %s", err)
		}
	}()

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced shutdown: %s", err)
	}
}

func buildHandler(r chi.Router, db *db.DB, logger logger.Logger, v validator.Validate, cfg *config.Config) chi.Router {
	// set up middleware
	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.Cors())

	authHandler := middleware.Auth(cfg.JwtSecret)
	// register user handlers
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo, v)
	userAPI := user.NewAPI(userService, v, authHandler)
	userAPI.RegisterHandlers(r)

	// register auth handlers
	authService := auth.NewService(userRepo, cfg.JwtSecret, cfg.JwtExpiration)
	authAPI := auth.NewAPI(authService, v)
	authAPI.RegisterHandlers(r)

	return r
}
