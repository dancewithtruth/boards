package user

import (
	"context"
	"net/http"
	"testing"

	"github.com/Wave-95/boards/backend-core/internal/endpoint"
	"github.com/Wave-95/boards/backend-core/internal/middleware"
	"github.com/Wave-95/boards/backend-core/internal/test"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestAPI(t *testing.T) {
	userRepo, jwtService, validator := getServiceDeps()
	userService := NewService(userRepo, validator)
	api := NewAPI(userService, jwtService, validator)
	router := chi.NewRouter()
	authHandler := middleware.Auth(jwtService)
	api.RegisterHandlers(router, authHandler)

	// Set up test user
	user := test.NewUser()
	userRepo.CreateUser(context.Background(), user)
	token, err := jwtService.GenerateToken(user.ID.String())
	if err != nil {
		assert.FailNow(t, "Failed to create test user", err)
	}
	header := test.AuthHeader(token)

	tests := []test.APITestCase{
		{
			Name:         "create user ok",
			Method:       http.MethodPost,
			URL:          "/users",
			Body:         `{"name":"John Doe", "is_guest": true}`,
			Header:       nil,
			WantStatus:   http.StatusCreated,
			WantResponse: "*token*",
		},
		{
			Name:         "create user no body",
			Method:       http.MethodPost,
			URL:          "/users",
			Body:         "",
			Header:       nil,
			WantStatus:   http.StatusBadRequest,
			WantResponse: `*` + endpoint.ErrMsgJSONDecode + `*`,
		},
		{
			Name:         "create user no name",
			Method:       http.MethodPost,
			URL:          "/users",
			Body:         `{"email": "test@email.com"}`,
			Header:       nil,
			WantStatus:   http.StatusBadRequest,
			WantResponse: "*name is a required field*",
		},
		{
			Name:         "get me ok",
			Method:       http.MethodGet,
			URL:          "/users/me",
			Body:         "",
			Header:       header,
			WantStatus:   http.StatusOK,
			WantResponse: `*` + user.ID.String() + `*`,
		},
	}

	for _, tc := range tests {
		test.Endpoint(t, router, tc)
	}

}
