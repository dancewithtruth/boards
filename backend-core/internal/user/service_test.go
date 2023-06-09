package user

import (
	"context"
	"testing"

	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/Wave-95/boards/backend-core/pkg/validator"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	validator := validator.New()
	mockRepo := &mockRepository{make(map[uuid.UUID]models.User)}
	service := NewService(mockRepo, validator)
	assert.NotNil(t, service)

	email := "testemail@gmail.com"
	password := "password123!"
	input := CreateUserInput{
		Name:     "Name",
		Email:    &email,
		Password: &password,
		IsGuest:  false,
	}
	t.Run("Create user", func(t *testing.T) {
		t.Run("with a valid user", func(t *testing.T) {
			user, err := service.CreateUser(context.Background(), input)
			assert.NoError(t, err)
			assert.Equal(t, input.Name, user.Name)
		})
	})
}
