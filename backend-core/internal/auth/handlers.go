package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Wave-95/boards/backend-core/internal/endpoint"
	"github.com/Wave-95/boards/backend-core/pkg/logger"
	"github.com/go-playground/validator/v10"
)

const (
	// ErrMsgBadLogin is a message displayed when a user provides bad login credentials
	ErrMsgBadLogin = "User does not exist"
	// ErrMsgInternalServer is a message displayed when an unexpected error occurs
	ErrMsgInternalServer = "Internal server error"
)

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
			endpoint.WriteWithError(w, http.StatusUnauthorized, ErrMsgBadLogin)
		case errors.As(err, &validator.ValidationErrors{}):
			endpoint.WriteValidationErr(w, input, err)
		default:
			logger.Errorf("handler: internal server error when performing login: %v", err)
			endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		}
		return
	}
	endpoint.WriteWithStatus(w, http.StatusOK, LoginDTO{Token: token})
}
