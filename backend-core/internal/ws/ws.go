package ws

import (
	"github.com/Wave-95/boards/backend-core/internal/board"
	"github.com/Wave-95/boards/backend-core/internal/jwt"
	"github.com/Wave-95/boards/backend-core/internal/post"
	"github.com/Wave-95/boards/backend-core/internal/user"
	"github.com/go-chi/chi/v5"
)

type WebSocket struct {
	userService  user.Service
	boardService board.Service
	postService  post.Service
	jwtService   jwt.Service
	boardHubs    map[string]*Hub
	destroy      chan string
}

func NewWebSocket(
	userService user.Service,
	boardService board.Service,
	postService post.Service,
	jwtService jwt.Service,
) *WebSocket {
	destroy := make(chan string)
	boardHubs := make(map[string]*Hub)
	go handleDestroy(destroy, boardHubs)
	return &WebSocket{
		userService:  userService,
		boardService: boardService,
		postService:  postService,
		jwtService:   jwtService,
		boardHubs:    boardHubs,
		destroy:      destroy,
	}
}

func (ws *WebSocket) RegisterHandlers(r chi.Router) {
	r.HandleFunc("/ws", ws.HandleConnection)
}
