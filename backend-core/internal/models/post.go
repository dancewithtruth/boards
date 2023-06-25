package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	// PostColorLightPink is a sample hex code background color for a post.
	PostColorLightPink = "#F5E6E8"
)

// Post defines the domain model for a Post entity.
type Post struct {
	ID        uuid.UUID `json:"id"`
	BoardID   uuid.UUID `json:"board_id"`
	UserID    uuid.UUID `json:"user_id"`
	Content   string    `json:"content"`
	PosX      int       `json:"pos_x"`
	PosY      int       `json:"pos_y"`
	Color     string    `json:"color"`
	Height    int       `json:"height"`
	ZIndex    int       `json:"z_index"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
