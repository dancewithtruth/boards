package post

import (
	"net/http"

	"github.com/Wave-95/boards/backend-core/internal/board"
	"github.com/Wave-95/boards/backend-core/pkg/validator"
	"github.com/go-chi/chi/v5"
)

type API struct {
	postService  Service
	boardService board.Service
	validator    validator.Validate
}

func NewAPI(postService Service, boardService board.Service, validator validator.Validate) API {
	return API{
		postService:  postService,
		boardService: boardService,
		validator:    validator,
	}
}

func (api *API) RegisterHandlers(r chi.Router, authHandler func(http.Handler) http.Handler) {
	r.Route("/posts", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(authHandler)
			r.Get("/", api.HandleListPosts)
		})
	})
}
