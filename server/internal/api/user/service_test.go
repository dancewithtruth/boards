package user

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	mockRepo := &mockRepository{make(map[uuid.UUID]User)}
	service := NewService(mockRepo)
	assert.NotNil(t, service)

	email := "testemail@gmail.com"
	password := "password123!"
	input := &CreateUserInput{
		Name:     "Name",
		Email:    &email,
		Password: &password,
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
			invalidEmail := "blah.com"

			invalidInput.Email = &invalidEmail
			_, err := service.CreateUser(invalidInput)
			assert.ErrorIs(t, err, ErrInvalidEmail)
		})

	})
}
