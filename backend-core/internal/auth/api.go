package auth

import (
	"github.com/Wave-95/boards/backend-core/pkg/validator"
	"github.com/go-chi/chi/v5"
)

// API encapsulates dependencies needed to perform auth-related duties.
type API struct {
	authService Service
	validator   validator.Validate
}

// NewAPI creates a new intance of the API struct.
func NewAPI(authService Service, validator validator.Validate) API {
	return API{
		authService: authService,
		validator:   validator,
	}
}

// RegisterHandlers registers the API's request handlers.
func (api *API) RegisterHandlers(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", api.HandleLogin)
	})
}
