package user

import (
	"net/http"

	"github.com/Wave-95/boards/backend-core/internal/jwt"
	"github.com/Wave-95/boards/backend-core/pkg/validator"
	"github.com/go-chi/chi/v5"
)

const (
	ErrMsgInternalServer = "Issue creating user"
)

type API struct {
	userService Service
	jwtService  jwt.Service
	validator   validator.Validate
}

func NewAPI(userService Service, jwtService jwt.Service, validator validator.Validate) API {
	return API{
		userService: userService,
		jwtService:  jwtService,
		validator:   validator,
	}
}

func (api *API) RegisterHandlers(r chi.Router, authHandler func(http.Handler) http.Handler) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", api.HandleCreateUser)
		r.Group(func(r chi.Router) {
			r.Use(authHandler)
			r.Get("/me", api.HandleGetUserMe)
		})
	})
}
