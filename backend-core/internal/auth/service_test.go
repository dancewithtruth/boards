package auth

import (
	"context"
	"testing"

	"github.com/Wave-95/boards/backend-core/internal/jwt"
	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/Wave-95/boards/backend-core/internal/test"
	"github.com/Wave-95/boards/backend-core/internal/user"
	"github.com/Wave-95/boards/backend-core/pkg/security"
	"github.com/Wave-95/boards/backend-core/pkg/validator"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	userRepo, jwtService, validate := getServiceDeps()
	authService := NewService(userRepo, jwtService, validate)
	assert.NotNil(t, authService)
	t.Run("Login", func(t *testing.T) {
		t.Run("using invalid input will return a validation error", func(t *testing.T) {
			invalidInput := LoginInput{
				Email:    "invalidemail",
				Password: "bad",
			}
			token, err := authService.Login(context.Background(), invalidInput)
			assert.True(t, validator.IsValidationError(err), "Expected a validation error", err)
			assert.Empty(t, token, "Expected an empty token to be returned")
		})

		t.Run("using credentials that exist will return a JWT token", func(t *testing.T) {
			user := setupUser(t, userRepo)
			defer cleanupUser(userRepo, user.ID)
			input := LoginInput{
				Email:    *user.Email,
				Password: *user.Password,
			}
			token, err := authService.Login(context.Background(), input)
			assert.NoError(t, err, "Expected no error when calling Login", err)
			assert.NotEmpty(t, token, "Expected a non-empty token to be returned")
		})

		t.Run("using credentials that do not exist will return a ErrBadLogin error", func(t *testing.T) {
			badInput := LoginInput{
				Email:    "bademail123@xyz.com",
				Password: "badpassword123",
			}
			token, err := authService.Login(context.Background(), badInput)
			assert.ErrorIs(t, err, errBadLogin, "Expected an ErrBadLogin error", err)
			assert.Empty(t, token, "Expected an empty token to be returned")
		})
	})
}

func getServiceDeps() (user.Repository, jwt.Service, validator.Validate) {
	userRepo := user.NewMockRepository()
	jwtService := test.NewJWTService()
	validate := validator.New()
	return userRepo, jwtService, validate
}

func setupUser(t *testing.T, userRepo user.Repository) models.User {
	user := test.NewUser()
	hashedPassword, err := security.HashPassword(*user.Password)
	if err != nil {
		assert.FailNow(t, "Failed to hash password.")
	}
	userWithHashedPW := user
	userWithHashedPW.Password = &hashedPassword
	err = userRepo.CreateUser(context.Background(), userWithHashedPW)
	if err != nil {
		assert.FailNow(t, "Failed to create test user.")
	}
	return user
}

func cleanupUser(userRepo user.Repository, userID uuid.UUID) {
	userRepo.DeleteUser(context.Background(), userID)
}
