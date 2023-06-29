package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	// PostColorLightPink is a sample hex code background color for a post.
	PostColorLightPink = "#F5E6E8"
)

// Post defines the domain model for a post entity.
type Post struct {
	ID          uuid.UUID `json:"id"`
	BoardID     uuid.UUID `json:"board_id"`
	UserID      uuid.UUID `json:"user_id"`
	Content     string    `json:"content"`
	Color       string    `json:"color"`
	Height      int       `json:"height"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	PostOrder   float64   `json:"post_order"`
	PostGroupID uuid.UUID `json:"post_group_id"`
}

// PostGroup defines the domain model for a post group entity.
type PostGroup struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	PosX      int       `json:"pos_x"`
	PosY      int       `json:"pos_y"`
	ZIndex    int       `json:"z_index"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
