package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID
	Name      string
	Email     *string
	Password  *string
	IsGuest   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
