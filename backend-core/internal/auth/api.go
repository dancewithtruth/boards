package auth

import (
	"github.com/Wave-95/boards/backend-core/pkg/validator"
	"github.com/go-chi/chi/v5"
)

type API struct {
	authService Service
	validator   validator.Validate
}

func NewAPI(authService Service, validator validator.Validate) API {
	return API{
		authService: authService,
		validator:   validator,
	}
}

func (api *API) RegisterHandlers(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", api.HandleLogin)
	})
}
