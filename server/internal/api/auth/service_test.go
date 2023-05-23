package auth

import (
	"context"
	"testing"
	"time"

	"github.com/Wave-95/boards/server/internal/api/user"
	"github.com/Wave-95/boards/server/internal/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	userRepo := user.NewMockRepository(make(map[uuid.UUID]*user.User))
	jwtService := jwt.New("abc123", 24)
	service := NewService(userRepo, jwtService)
	assert.NotNil(t, service)
	testUser := newTestUser()
	t.Run("Login", func(t *testing.T) {
		t.Run("with valid credentials returns auth token", func(t *testing.T) {
			ctx := context.Background()
			userRepo.CreateUser(ctx, testUser)
			token, err := service.Login(ctx, *testUser.Email, *testUser.Password)
			assert.NotEmpty(t, token, "expected a non-empty auth token")
			assert.NoError(t, err, "expected no error when calling Login")
		})

		t.Run("with bad credentials returns error", func(t *testing.T) {
			ctx := context.Background()
			userRepo.CreateUser(ctx, testUser)
			badEmail := "bademail123.com"
			badPassword := "badpassword123"
			token, err := service.Login(ctx, badEmail, badPassword)
			assert.Empty(t, token, "expected an empty auth token")
			assert.ErrorIs(t, err, ErrBadLogin, "expected an ErrUserDoesNotExist when calling Login")
		})
	})
}

func newTestUser() *user.User {
	email := "johndoe@gmail.com"
	password := "password123"
	user := &user.User{
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
