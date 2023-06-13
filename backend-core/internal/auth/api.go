package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Wave-95/boards/backend-core/internal/endpoint"
	"github.com/Wave-95/boards/backend-core/pkg/logger"
	"github.com/Wave-95/boards/backend-core/pkg/validator"
	"github.com/go-chi/chi/v5"
	v "github.com/go-playground/validator/v10"
)

const (
	// ErrMsgUserDoesNotExist is a message displayed when a user provides login credentials for a user that does not exist
	ErrMsgUserDoesNotExist = "User does not exist."
	// ErrMsgInternalServer is a message displayed when an unexpected error occurs
	ErrMsgInternalServer = "Internal server error."
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

// HandleLogin handles a user's login request. It returns a token in the response
// if the login is successful.
func (api *API) HandleLogin(w http.ResponseWriter, r *http.Request) {
	logger := logger.FromContext(r.Context())
	var input LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		endpoint.HandleDecodeErr(w, err)
		return
	}
	defer r.Body.Close()
	token, err := api.authService.Login(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, ErrBadLogin):
			endpoint.WriteWithError(w, http.StatusNotFound, ErrMsgUserDoesNotExist)
		case errors.As(err, &v.ValidationErrors{}):
			endpoint.WriteValidationErr(w, input, err)
		default:
			logger.Errorf("HandleLogin: Failed to login user due to internal server error > %v", err)
			endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		}
		return
	}
	endpoint.WriteWithStatus(w, http.StatusOK, LoginDTO{Token: token})
}

// RegisterHandlers registers the API's request handlers.
func (api *API) RegisterHandlers(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", api.HandleLogin)
	})
}
