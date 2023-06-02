package board

import (
	"net/http"

	"github.com/Wave-95/boards/server/pkg/validator"
	"github.com/go-chi/chi/v5"
)

type API struct {
	boardService Service
	validator    validator.Validate
}

func NewAPI(boardService Service, validator validator.Validate) API {
	return API{
		boardService: boardService,
		validator:    validator,
	}
}

func (api *API) RegisterHandlers(r chi.Router, authHandler func(http.Handler) http.Handler) {
	r.Route("/boards", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(authHandler)
			r.Get("/", api.HandleGetBoards)
			r.Get("/{boardId}", api.HandleGetBoard)
			r.Post("/", api.HandleCreateBoard)
		})
	})
}
