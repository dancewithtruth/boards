package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	PostColorLightPink = "#F5E6E8"
)

type Post struct {
	Id        uuid.UUID `json:"id"`
	BoardId   uuid.UUID `json:"board_id"`
	UserId    uuid.UUID `json:"user_id"`
	Content   string    `json:"content"`
	PosX      int       `json:"pos_x"`
	PosY      int       `json:"pos_y"`
	Color     string    `json:"color"`
	Height    float64   `json:"height"`
	ZIndex    int       `json:"z_index"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
