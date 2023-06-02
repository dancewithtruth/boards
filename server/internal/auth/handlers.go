package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Wave-95/boards/server/internal/endpoint"
)

const (
	ErrMsgBadLogin       = "User does not exist"
	ErrMsgInternalServer = "Issue logging in"
)

func (api *API) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var input LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		endpoint.HandleDecodeErr(w, err)
		return
	}
	defer r.Body.Close()

	if err := api.validator.Struct(input); err != nil {
		endpoint.WriteValidationErr(w, input, err)
		return
	}

	token, err := api.authService.Login(r.Context(), input)
	if err != nil {
		if errors.Is(err, ErrBadLogin) {
			endpoint.WriteWithError(w, http.StatusUnauthorized, ErrMsgBadLogin)
			return
		}
		endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		return
	}
	endpoint.WriteWithStatus(w, http.StatusOK, LoginDTO{Token: token})
}
