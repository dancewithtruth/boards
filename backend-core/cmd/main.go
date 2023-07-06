package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Wave-95/boards/backend-core/db"
	"github.com/Wave-95/boards/backend-core/internal/auth"
	"github.com/Wave-95/boards/backend-core/internal/board"
	"github.com/Wave-95/boards/backend-core/internal/config"
	"github.com/Wave-95/boards/backend-core/internal/endpoint"
	"github.com/Wave-95/boards/backend-core/internal/jwt"
	"github.com/Wave-95/boards/backend-core/internal/middleware"
	"github.com/Wave-95/boards/backend-core/internal/post"
	"github.com/Wave-95/boards/backend-core/internal/user"
	"github.com/Wave-95/boards/backend-core/internal/ws"
	"github.com/Wave-95/boards/backend-core/pkg/logger"
	"github.com/Wave-95/boards/backend-core/pkg/validator"
	"github.com/go-chi/chi/v5"
)

func main() {
	logger := logger.New()
	validator := validator.New()

	// Get config
	cfg, err := config.Load(".env")
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}

	// Connect to database
	conn, err := db.Connect(cfg.DatabaseConfig)
	if err != nil {
		logger.Fatalf("Error connecting to db: %v", err)
	}
	defer conn.Close()

	// Setup server
	r := chi.NewRouter()
	server := http.Server{Addr: cfg.ServerPort, Handler: buildHandler(r, conn, logger, validator, cfg)}

	// Graceful shutdown
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

// buildHandler sets up all the middleware and API routes for the server.
func buildHandler(r chi.Router, db *db.DB, logger logger.Logger, v validator.Validate, cfg *config.Config) chi.Router {
	// Set up middleware
	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.Cors())

	// Set up repositories
	userRepo := user.NewRepository(db)
	boardRepo := board.NewRepository(db)
	postRepo := post.NewRepository(db)

	// Set up services
	jwtService := jwt.New(cfg.JwtSecret, cfg.JwtExpiration)
	authService := auth.NewService(userRepo, jwtService, v)
	userService := user.NewService(userRepo, v)
	boardService := board.NewService(boardRepo, v)
	postService := post.NewService(postRepo)

	// Set up APIs
	userAPI := user.NewAPI(userService, jwtService, v)
	authAPI := auth.NewAPI(authService, v)
	boardAPI := board.NewAPI(boardService, v)
	postAPI := post.NewAPI(postService, boardService, v)
	websocket := ws.NewWebSocket(userService, boardService, postService, jwtService)

	// Set up auth handler
	authHandler := middleware.Auth(jwtService)

	// Register handlers
	userAPI.RegisterHandlers(r, authHandler)
	authAPI.RegisterHandlers(r)
	boardAPI.RegisterHandlers(r, authHandler)
	postAPI.RegisterHandlers(r, authHandler)
	websocket.RegisterHandlers(r)
	r.Get("/ping", handlePingCheck)

	return r
}

func handlePingCheck(w http.ResponseWriter, r *http.Request) {
	endpoint.WriteWithStatus(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{Message: "pong"})
}
