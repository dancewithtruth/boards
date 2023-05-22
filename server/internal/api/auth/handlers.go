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

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (api *API) HandleLogin(w http.ResponseWriter, r *http.Request) {
	req := LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		endpoint.HandleDecodeErr(w, err)
		return
	}
	defer r.Body.Close()

	if err := api.validator.Struct(req); err != nil {
		endpoint.HandleValidationErr(w, err)
		return
	}

	token, err := api.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, ErrBadLogin) {
			endpoint.WriteWithError(w, http.StatusUnauthorized, ErrMsgBadLogin)
			return
		}
		endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		return
	}
	endpoint.WriteWithStatus(w, http.StatusOK, LoginResponse{Token: token})
}
