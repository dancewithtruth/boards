package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	mockRepo := &mockRepository{}
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
			err := service.CreateUser(input)
			assert.NoError(t, err)
		})

		t.Run("with an invalid email", func(t *testing.T) {
			invalidInput := input
			invalidInput.Email = "xyz.com"
			err := service.CreateUser(invalidInput)
			assert.ErrorIs(t, err, ErrInvalidEmail)
		})

	})
}
