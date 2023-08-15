package ws

import (
	"github.com/Wave-95/boards/backend-core/internal/board"
	"github.com/Wave-95/boards/backend-core/internal/jwt"
	"github.com/Wave-95/boards/backend-core/internal/post"
	"github.com/Wave-95/boards/backend-core/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

type WebSocket struct {
	userService  user.Service
	boardService board.Service
	postService  post.Service
	jwtService   jwt.Service
	rdb          *redis.Client
}

func NewWebSocket(
	userService user.Service,
	boardService board.Service,
	postService post.Service,
	jwtService jwt.Service,
	rdb *redis.Client,
) *WebSocket {

	return &WebSocket{
		userService:  userService,
		boardService: boardService,
		postService:  postService,
		jwtService:   jwtService,
		rdb:          rdb,
	}
}

func (ws *WebSocket) RegisterHandlers(r chi.Router) {
	r.HandleFunc("/ws", ws.HandleConnection)
}
