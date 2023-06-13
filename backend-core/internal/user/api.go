package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Wave-95/boards/backend-core/internal/endpoint"
	"github.com/Wave-95/boards/backend-core/internal/jwt"
	"github.com/Wave-95/boards/backend-core/internal/middleware"
	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/Wave-95/boards/backend-core/pkg/logger"
	"github.com/Wave-95/boards/backend-core/pkg/validator"
	"github.com/go-chi/chi/v5"
	v "github.com/go-playground/validator/v10"
)

const (
	// ErrMsgInternalServer is an error message for unexpected errors
	ErrMsgInternalServer = "Internal server error"
	// ErrMsgInvalidSearchParam is an error message for an invalid search query parameter
	ErrMsgInvalidSearchParam = `Invalid or missing search param. Try using "email".`
)

// API represents a struct for the user API
type API struct {
	userService Service
	jwtService  jwt.Service
	validator   validator.Validate
}

// NewAPI initializes an API struct
func NewAPI(userService Service, jwtService jwt.Service, validator validator.Validate) API {
	return API{
		userService: userService,
		jwtService:  jwtService,
		validator:   validator,
	}
}

// HandleCreateUser creates a user and generates a JWT token using the user ID as a field.
func (api *API) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logger.FromContext(ctx)

	// Decode request
	var input CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Errorf("handler: failed to decode request: %v", err)
		endpoint.HandleDecodeErr(w, err)
		return
	}
	defer r.Body.Close()

	// Validate request
	if err := api.validator.Struct(input); err != nil {
		logger.Errorf("handler: failed to validate request: %v", err)
		endpoint.WriteValidationErr(w, input, err)
		return
	}

	// Create user and handle errors
	user, err := api.userService.CreateUser(ctx, input)
	if err != nil {
		switch {
		case errors.As(err, &v.ValidationErrors{}):
			endpoint.WriteWithError(w, http.StatusBadRequest, validator.GetValidationErrMsg(input, err))
		case errors.Is(err, ErrEmailAlreadyExists):
			endpoint.WriteWithError(w, http.StatusConflict, ErrEmailAlreadyExists.Error())
		default:
			endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		}
		return
	}

	jwtToken, err := api.jwtService.GenerateToken(user.ID.String())
	if err != nil {
		endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
	}
	endpoint.WriteWithStatus(w, http.StatusCreated, CreateUserDTO{User: user, JwtToken: jwtToken})
}

// HandleGetUserMe is protected with an authHandler and expects the userID to be present
// on the request context
func (api *API) HandleGetUserMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.UserIDFromContext(ctx)
	user, err := api.userService.GetUser(ctx, userID)
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
	logger := logger.FromContext(ctx)
	queryParams := r.URL.Query()
	email := queryParams.Get("email")
	if email == "" {
		endpoint.WriteWithError(w, http.StatusBadRequest, ErrMsgInvalidSearchParam)
		return
	}
	users, err := api.userService.ListUsersByFuzzyEmail(ctx, email)
	if err != nil {
		logger.Errorf("handler: failed to list users by search: %w", err)
		switch {
		default:
			endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		}
		return
	}
	endpoint.WriteWithStatus(w, http.StatusOK, struct {
		Result []models.User `json:"result"`
	}{Result: users})
}

// RegisterHandlers is a function that registers all the handlers for the user endpoints
func (api *API) RegisterHandlers(r chi.Router, authHandler func(http.Handler) http.Handler) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", api.HandleCreateUser)
		r.Get("/search", api.HandleListUsersBySearch)
		r.Group(func(r chi.Router) {
			r.Use(authHandler)
			r.Get("/me", api.HandleGetUserMe)
		})
	})
}
