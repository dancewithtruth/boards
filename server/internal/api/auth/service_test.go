package auth

import (
	"context"
	"testing"
	"time"

	"github.com/Wave-95/boards/server/internal/api/user"
	"github.com/Wave-95/boards/server/internal/jwt"
	"github.com/Wave-95/boards/server/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	userRepo := user.NewMockRepository(make(map[uuid.UUID]models.User))
	jwtService := jwt.New("abc123", 24)
	service := NewService(userRepo, jwtService)
	assert.NotNil(t, service)
	testUser := newTestUser()
	input := LoginInput{
		Email:    *testUser.Email,
		Password: *testUser.Password,
	}
	t.Run("Login", func(t *testing.T) {
		t.Run("with valid credentials returns auth token", func(t *testing.T) {
			ctx := context.Background()
			userRepo.CreateUser(ctx, testUser)
			token, err := service.Login(ctx, input)
			assert.NotEmpty(t, token, "expected a non-empty auth token")
			assert.NoError(t, err, "expected no error when calling Login")
		})

		t.Run("with bad credentials returns error", func(t *testing.T) {
			ctx := context.Background()
			userRepo.CreateUser(ctx, testUser)
			badInput := LoginInput{
				Email:    "bademail123.com",
				Password: "badpassword123",
			}
			token, err := service.Login(ctx, badInput)
			assert.Empty(t, token, "expected an empty auth token")
			assert.ErrorIs(t, err, ErrBadLogin, "expected an ErrUserDoesNotExist when calling Login")
		})
	})
}

func newTestUser() models.User {
	email := "johndoe@gmail.com"
	password := "password123"
	user := models.User{
		Id:        uuid.New(),
		Name:      "testname",
		Email:     &email,
		Password:  &password,
		IsGuest:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return user
}
