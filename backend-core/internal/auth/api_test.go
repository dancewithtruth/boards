package auth

import (
	"net/http"
	"testing"

	"github.com/Wave-95/boards/backend-core/internal/test"
	"github.com/go-chi/chi/v5"
)

func TestAPI(t *testing.T) {
	userRepo, jwtService, validator := getServiceDeps()
	authService := NewService(userRepo, jwtService, validator)
	api := NewAPI(authService, validator)
	router := chi.NewRouter()
	api.RegisterHandlers(router)
	url := "/auth/login"

	user := setupUser(t, userRepo)
	defer cleanupUser(userRepo, user.ID)

	tests := []test.APITestCase{
		{
			Name:         "invalid login returns error response",
			Method:       http.MethodPost,
			URL:          url,
			Body:         `{"email":"invalidemail", "password": "password"}`,
			Header:       nil,
			WantStatus:   http.StatusBadRequest,
			WantResponse: "*Invalid input on email*",
		},
		{
			Name:         "ok login returns token",
			Method:       http.MethodPost,
			URL:          url,
			Body:         `{"email":"` + *user.Email + `", "password": "` + *user.Password + `"}`,
			Header:       nil,
			WantStatus:   http.StatusOK,
			WantResponse: "*token*",
		},
		{
			Name:         "not found login returns error response",
			Method:       http.MethodPost,
			URL:          url,
			Body:         `{"email":"notfound@example.com", "password": "notfound"}`,
			Header:       nil,
			WantStatus:   http.StatusNotFound,
			WantResponse: `*` + ErrMsgUserDoesNotExist + `*`,
		},
	}

	for _, tc := range tests {
		test.Endpoint(t, router, tc)
	}
}
