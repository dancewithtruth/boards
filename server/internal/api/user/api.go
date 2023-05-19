package user

import (
	"errors"
	"net/http"
	"time"

	"github.com/Wave-95/boards/server/internal/endpoint"
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
	Name     string  `json:"name" validate:"required"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
	IsGuest  bool    `json:"is_guest" validate:"omitempty,required"`
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
	Email     *string   `json:"email"`
	IsGuest   bool      `json:"is_guest"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (api *API) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	// decode request and validate
	var createUserRequest CreateUserRequest
	err := endpoint.DecodeAndValidate(w, r, &createUserRequest, ErrInvalidCreateUserRequest)
	if err != nil {
		endpoint.WriteWithError(w, http.StatusBadRequest, err)
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
		endpoint.WriteWithError(w, http.StatusInternalServerError, ErrInternalServerError)
		return
	}

	// write response
	createUserResponse := CreateUserResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		IsGuest:   user.IsGuest,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	endpoint.WriteWithStatus(w, http.StatusCreated, createUserResponse)
}
