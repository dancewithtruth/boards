package ws

import (
	"fmt"

	"github.com/Wave-95/boards/backend-core/internal/board"
	"github.com/Wave-95/boards/backend-core/internal/config"
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
	boardHubs    map[string]*Hub
	destroy      chan string
}

func NewWebSocket(
	userService user.Service,
	boardService board.Service,
	postService post.Service,
	jwtService jwt.Service,
	rdbConfig config.RedisConfig,
) *WebSocket {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%v:%v", rdbConfig.Host, rdbConfig.Port),
	})

	boardHubs := make(map[string]*Hub)
	destroy := make(chan string)

	go handleDestroy(destroy, boardHubs)

	return &WebSocket{
		userService:  userService,
		boardService: boardService,
		postService:  postService,
		jwtService:   jwtService,
		rdb:          rdb,
		boardHubs:    boardHubs,
		destroy:      destroy,
	}
}

func (ws *WebSocket) RegisterHandlers(r chi.Router) {
	r.HandleFunc("/ws", ws.HandleConnection)
}
