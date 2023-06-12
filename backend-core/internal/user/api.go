package user

import (
	"net/http"

	"github.com/Wave-95/boards/backend-core/internal/jwt"
	"github.com/Wave-95/boards/backend-core/pkg/validator"
	"github.com/go-chi/chi/v5"
)

const (
	// ErrMsgInternalServer is an error message for unexpected errors
	ErrMsgInternalServer = "Internal server error"
	// ErrMsgInvalidSearchParam is an error message for an invalid search query parameter
	ErrMsgInvalidSearchParam = `Invalid or missing search param. Try using "email".`
)

// API represents a struct for the user API
type API struct {
	userService Service
	jwtService  jwt.Service
	validator   validator.Validate
}

// NewAPI initializes an API struct
func NewAPI(userService Service, jwtService jwt.Service, validator validator.Validate) API {
	return API{
		userService: userService,
		jwtService:  jwtService,
		validator:   validator,
	}
}

// RegisterHandlers is a function that registers all the handlers for the user endpoints
func (api *API) RegisterHandlers(r chi.Router, authHandler func(http.Handler) http.Handler) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", api.HandleCreateUser)
		r.Get("/search", api.HandleListUsersBySearch)
		r.Group(func(r chi.Router) {
			r.Use(authHandler)
			r.Get("/me", api.HandleGetUserMe)
		})
	})
}
