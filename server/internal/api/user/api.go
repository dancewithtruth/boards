package user

import (
	"github.com/Wave-95/boards/server/pkg/validator"
	"github.com/go-chi/chi/v5"
)

const (
	ErrMsgInternalServer = "Issue creating user"
)

type API struct {
	service   Service
	validator validator.Validate
}

func NewAPI(service Service, validator validator.Validate) API {
	return API{service: service, validator: validator}
}

func (api *API) RegisterHandlers(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", api.HandleCreateUser)
	})
}
