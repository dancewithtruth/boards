package test

import (
	"time"

	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/google/uuid"
)

type NewUserOpt func(user *models.User)

// WithEmail is an optional closure function to be passed into NewUser
func WithEmail(email string) NewUserOpt {
	return func(user *models.User) {
		user.Email = &email
	}
}

// NewUser generates a test user model
func NewUser(opts ...NewUserOpt) models.User {
	email := uuid.New().String() + "@example.com"
	password := "password123"
	user := models.User{
		ID:        uuid.New(),
		Name:      "testname",
		Email:     &email,
		Password:  &password,
		IsGuest:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	for _, opt := range opts {
		opt(&user)
	}
	return user
}
