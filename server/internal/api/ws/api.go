package ws

import (
	"github.com/Wave-95/boards/server/internal/jwt"
	"github.com/Wave-95/boards/server/pkg/validator"
	"github.com/go-chi/chi/v5"
)

type API struct {
	hub        *Hub
	jwtService jwt.JWTService
	validator  validator.Validate
}

func NewAPI(jwtService jwt.JWTService, validator validator.Validate) API {
	hub := newHub()
	go hub.run()
	return API{
		hub:        hub,
		jwtService: jwtService,
		validator:  validator,
	}
}

func (api *API) RegisterHandlers(r chi.Router) {
	r.HandleFunc("/ws", api.HandleWebSocket)
}
