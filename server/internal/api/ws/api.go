package ws

import (
	"github.com/Wave-95/boards/server/internal/jwt"
	"github.com/go-chi/chi/v5"
)

type API struct {
	hub        *Hub
	jwtService jwt.JWTService
}

func NewAPI(jwtService jwt.JWTService) API {
	hub := newHub()
	go hub.run()
	return API{hub: hub, jwtService: jwtService}
}

func (api *API) RegisterHandlers(r chi.Router) {
	r.HandleFunc("/ws", api.HandleWebSocket)
}
