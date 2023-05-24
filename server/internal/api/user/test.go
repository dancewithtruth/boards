package user

import (
	"time"

	"github.com/google/uuid"
)

func NewTestUser() *User {
	email := uuid.New().String() + "email.com"
	password := "password123"
	user := &User{
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
