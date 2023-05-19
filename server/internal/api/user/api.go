package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Wave-95/boards/server/internal/response"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var (
	ErrMissingName              = errors.New("Missing name")
	ErrInvalidCreateUserRequest = errors.New("Invalid create user request")
	ErrInternalServerError      = errors.New("Issue creating user")
)

type API struct {
	service Service
}

func NewAPI(service Service) API {
	return API{service: service}
}

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsGuest  bool   `json:"is_guest"`
}

func (req *CreateUserRequest) Validate() error {
	v := validator.New()
	if err := v.Struct(req); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Name" {
				return ErrMissingName
			}
		}
		return ErrInvalidCreateUserRequest
	}
	return nil
}

type CreateUserResponse struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	IsGuest   bool      `json:"is_guest"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (api *API) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	// decode request and validate
	var createUserRequest CreateUserRequest
	json.NewDecoder(r.Body).Decode(&createUserRequest)
	if err := createUserRequest.Validate(); err != nil {
		response.WriteWithError(w, http.StatusBadRequest, err)
		return
	}

	// create user
	input := CreateUserInput{
		Name:     createUserRequest.Name,
		Email:    createUserRequest.Email,
		Password: createUserRequest.Password,
		IsGuest:  createUserRequest.IsGuest,
	}
	user, err := api.service.CreateUser(input)
	if err != nil {
		response.WriteWithError(w, http.StatusInternalServerError, ErrInternalServerError)
		return
	}

	// write response
	w.WriteHeader(http.StatusCreated)
	createUserResponse := CreateUserResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		IsGuest:   user.IsGuest,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	response.WriteWithStatus(w, http.StatusCreated, createUserResponse)
}
