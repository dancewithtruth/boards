package user

import "github.com/go-chi/chi/v5"

func (api *API) RegisterHandlers(r chi.Router) {
	r.Post("/users", api.HandleCreateUser)
}
