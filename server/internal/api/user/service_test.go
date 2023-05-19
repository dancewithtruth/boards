package user

import (
	"testing"

	"github.com/Wave-95/boards/server/internal/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	mockRepo := &mockRepository{make(map[uuid.UUID]entity.User)}
	service := NewService(mockRepo)
	assert.NotNil(t, service)

	input := CreateUserInput{
		Name:     "Name",
		Email:    "testemail@gmail.com",
		Password: "password123!",
		IsGuest:  false,
	}
	t.Run("Create user", func(t *testing.T) {
		t.Run("with a valid user", func(t *testing.T) {
			user, err := service.CreateUser(input)
			assert.NoError(t, err)
			assert.Equal(t, input.Name, user.Name)
		})

		t.Run("with an invalid email", func(t *testing.T) {
			invalidInput := input
			invalidInput.Email = "xyz.com"
			_, err := service.CreateUser(invalidInput)
			assert.ErrorIs(t, err, ErrInvalidEmail)
		})

	})
}
