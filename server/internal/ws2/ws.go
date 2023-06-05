package ws2

import (
	"github.com/Wave-95/boards/server/internal/board"
	"github.com/Wave-95/boards/server/internal/jwt"
	"github.com/Wave-95/boards/server/internal/user"
	"github.com/go-chi/chi/v5"
)

type WebSocket struct {
	userService  user.Service
	boardService board.Service
	jwtService   jwt.Service
	boardHubs    map[string]*Hub
	destroy      chan string
}

func NewWebSocket(userService user.Service, boardService board.Service, jwtService jwt.Service) *WebSocket {
	destroy := make(chan string)
	boardHubs := make(map[string]*Hub)
	go handleDestroy(destroy, boardHubs)
	return &WebSocket{
		userService:  userService,
		boardService: boardService,
		jwtService:   jwtService,
		boardHubs:    boardHubs,
		destroy:      destroy,
	}
}

func (ws *WebSocket) RegisterHandlers(r chi.Router) {
	r.HandleFunc("/ws", ws.HandleConnection)
}
