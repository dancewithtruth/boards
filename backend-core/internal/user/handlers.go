package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Wave-95/boards/backend-core/internal/endpoint"
	"github.com/Wave-95/boards/backend-core/internal/middleware"
	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/Wave-95/boards/backend-core/pkg/logger"
	"github.com/Wave-95/boards/backend-core/pkg/validator"
	v "github.com/go-playground/validator/v10"
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
		endpoint.WriteValidationErr(w, input, err)
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
	endpoint.WriteWithStatus(w, http.StatusCreated, CreateUserDTO{User: user, JwtToken: jwtToken})
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

// HandleListUsersBySearch looks for an email query param and lists the top 10 closest user matches
// by email input
func (api *API) HandleListUsersBySearch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	queryParams := r.URL.Query()
	email := queryParams.Get("email")
	if email == "" {
		endpoint.WriteWithError(w, http.StatusBadRequest, ErrMsgInvalidSearchParam)
		return
	}
	input := ListUsersByFuzzyEmailInput{Email: email}
	users, err := api.userService.ListUsersByFuzzyEmail(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, v.ValidationErrors{}):
			endpoint.WriteWithError(w, http.StatusConflict, validator.GetValidationErrMsg(input, err))
		default:
			endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		}
		return
	}
	endpoint.WriteWithStatus(w, http.StatusOK, struct {
		Result []models.User `json:"result"`
	}{Result: users})
}
