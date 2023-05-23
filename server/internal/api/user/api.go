package user

import (
	"net/http"

	"github.com/Wave-95/boards/server/pkg/validator"
	"github.com/go-chi/chi/v5"
)

const (
	ErrMsgInternalServer = "Issue creating user"
)

type API struct {
	service     Service
	validator   validator.Validate
	authHandler func(http.Handler) http.Handler
}

func NewAPI(service Service, validator validator.Validate, authHandler func(http.Handler) http.Handler) API {
	return API{
		service:     service,
		validator:   validator,
		authHandler: authHandler,
	}
}

func (api *API) RegisterHandlers(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", api.HandleCreateUser)
		r.Group(func(r chi.Router) {
			r.Use(api.authHandler)
			r.Get("/me", api.HandleGetUserMe)
		})
	})
}
