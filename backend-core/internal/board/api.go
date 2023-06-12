package board

import (
	"net/http"

	"github.com/Wave-95/boards/backend-core/pkg/validator"
	"github.com/go-chi/chi/v5"
)

// API encapsulates dependencies needed to perform board related duties.
type API struct {
	boardService Service
	validator    validator.Validate
}

// NewAPI creates a new intance of the API struct.
func NewAPI(boardService Service, validator validator.Validate) API {
	return API{
		boardService: boardService,
		validator:    validator,
	}
}

// RegisterHandlers registers the API's request handlers.
func (api *API) RegisterHandlers(r chi.Router, authHandler func(http.Handler) http.Handler) {
	r.Route("/boards", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(authHandler)
			r.Get("/", api.HandleGetBoards)
			r.Post("/", api.HandleCreateBoard)

			r.Route("/{boardID}", func(r chi.Router) {
				r.Get("/", api.HandleGetBoard)
				r.Post("/invites", api.HandleCreateBoardInvites)
			})
		})
	})
}
