package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Wave-95/boards/server/internal/endpoint"
	"github.com/Wave-95/boards/server/pkg/logger"
	"github.com/Wave-95/boards/server/pkg/validator"
)

var (
	ErrMissingName              = errors.New("Missing name")
	ErrInvalidCreateUserRequest = errors.New("Invalid create user request")
	ErrInternalServerError      = errors.New("Issue creating user")
)

type API struct {
	service   Service
	validator validator.Validate
}

func NewAPI(service Service, validator validator.Validate) API {
	return API{service: service, validator: validator}
}

func (api *API) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logger.FromContext(ctx)

	// decode request
	var request CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Errorf("Issue decoding request: %v", err)
		err = endpoint.HandleDecodeErr(err, ErrInvalidCreateUserRequest)
		endpoint.WriteWithError(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()

	// validate request
	if err := api.validator.Struct(request); err != nil {
		logger.Errorf("Issue validating request: %v", err)
		endpoint.WriteWithError(w, http.StatusBadRequest, err)
	}

	// create user
	user, err := api.service.CreateUser(ctx, request.ToInput())
	if err != nil {
		endpoint.WriteWithError(w, http.StatusInternalServerError, ErrInternalServerError)
		return
	}

	// write response
	endpoint.WriteWithStatus(w, http.StatusCreated, user.ToDto())
}
