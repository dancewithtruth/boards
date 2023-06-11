package test

import (
	"time"

	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/google/uuid"
)

type NewUserOpt func(user *models.User)

func WithEmail(email string) NewUserOpt {
	return func(user *models.User) {
		user.Email = &email
	}
}

func NewUser(opts ...NewUserOpt) models.User {
	email := uuid.New().String() + "email.com"
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
	for _, opt := range opts {
		opt(&user)
	}
	return user
}
