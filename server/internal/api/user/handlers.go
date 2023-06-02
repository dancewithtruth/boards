package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Wave-95/boards/server/internal/endpoint"
	"github.com/Wave-95/boards/server/internal/middleware"
	"github.com/Wave-95/boards/server/internal/models"
	"github.com/Wave-95/boards/server/pkg/logger"
)

func (api *API) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logger.FromContext(ctx)

	// decode request
	var input CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Errorf("handler: failed to decode request: %v", err)
		endpoint.HandleDecodeErr(w, err)
		return
	}
	defer r.Body.Close()

	// validate request
	if err := api.validator.Struct(input); err != nil {
		logger.Errorf("handler: failed to validate request: %v", err)
		endpoint.WriteValidationErr(w, err)
		return
	}

	// create user and handle errors
	user, err := api.userService.CreateUser(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmailAlreadyExists):
			endpoint.WriteWithError(w, http.StatusConflict, ErrEmailAlreadyExists.Error())
		default:
			endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		}
		return
	}

	jwtToken, err := api.jwtService.GenerateToken(user.Id.String())
	if err != nil {
		endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
	}

	// write response
	endpoint.WriteWithStatus(w, http.StatusCreated, struct {
		User     models.User `json:"user"`
		JwtToken string      `json:"jwt_token"`
	}{user, jwtToken})
}

// HandleGetUserMe is protected with an authHandler and expects the userID to be present
// on the request context
func (api *API) HandleGetUserMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := middleware.UserIdFromContext(ctx)
	user, err := api.userService.GetUser(ctx, userId)
	if err != nil {
		switch {
		default:
			endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		}
		return
	}
	// write response
	endpoint.WriteWithStatus(w, http.StatusOK, user)
}

// Inputs
type CreateUserInput struct {
	Name     string  `json:"name" validate:"required,min=2,max=12"`
	Email    *string `json:"email" validate:"omitempty,email,required"`
	Password *string `json:"password" validate:"omitempty,min=8"`
	IsGuest  bool    `json:"is_guest" validate:"omitempty,required"`
}
