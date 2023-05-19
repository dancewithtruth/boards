package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Wave-95/boards/server/internal/endpoint"
	"github.com/Wave-95/boards/server/pkg/logger"
)

func (api *API) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logger.FromContext(ctx)

	// decode request
	var request CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Errorf("handler: failed to decode request: %v", err)
		endpoint.HandleDecodeErr(w, err)
		return
	}
	defer r.Body.Close()

	// validate request
	if err := api.validator.Struct(request); err != nil {
		logger.Errorf("handler: failed to validate request: %v", err)
		endpoint.HandleValidationErr(w, err)
		return
	}

	// create user and handle errors
	user, err := api.service.CreateUser(ctx, request.ToInput())
	if err != nil {
		switch {
		case errors.Is(err, ErrEmailAlreadyExists):
			endpoint.WriteWithError(w, http.StatusConflict, ErrEmailAlreadyExists.Error())
		default:
			endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		}
		return
	}

	// write response
	endpoint.WriteWithStatus(w, http.StatusCreated, user.ToDto())
}
