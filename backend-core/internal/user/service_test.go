package user

import (
	"context"
	"testing"

	"github.com/Wave-95/boards/backend-core/internal/jwt"
	"github.com/Wave-95/boards/backend-core/internal/test"
	"github.com/Wave-95/boards/backend-core/pkg/validator"
	"github.com/Wave-95/boards/wrappers/amqp"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	userRepo, _, validate := getServiceDeps()
	amqp := amqp.NewMock()
	userService := NewService(userRepo, amqp, validate)
	assert.NotNil(t, userService)
	t.Run("Create and get user", func(t *testing.T) {
		t.Run("using valid user input", func(t *testing.T) {
			// Create
			input := CreateUserInput{
				Name:    "john doe",
				IsGuest: true,
			}
			newUser, err := userService.CreateUser(context.Background(), input)
			assert.NoError(t, err)
			assert.Equal(t, "John Doe", newUser.Name, "Expected first letter of each word in name to be capitalized.")
			// Get
			user, err := userService.GetUser(context.Background(), newUser.ID.String())
			assert.NoError(t, err)
			assert.Equal(t, newUser.ID, user.ID)
		})

		t.Run("using invalid user input", func(t *testing.T) {
			input := CreateUserInput{}
			_, err := userService.CreateUser(context.Background(), input)
			assert.True(t, validator.IsValidationError(err), "Expected error to be a validation error")
		})
	})
}

func getServiceDeps() (Repository, jwt.Service, validator.Validate) {
	userRepo := NewMockRepository()
	jwtService := test.NewJWTService()
	validator := validator.New()
	return userRepo, jwtService, validator
}
