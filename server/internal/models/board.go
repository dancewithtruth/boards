package models

import (
	"time"

	"github.com/google/uuid"
)

type Board struct {
	Id          uuid.UUID `json:"id"`
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
	UserId      uuid.UUID `json:"user_id"`
	Users       []User    `json:"users"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (b *Board) ToDto() BoardResponse {
	return BoardResponse{
		Id:          b.Id,
		Name:        b.Name,
		Description: b.Description,
		UserId:      b.UserId,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
}

type CreateBoardRequest struct {
	Name        *string `json:"name" validate:"omitempty,required,min=3,max=20"`
	Description *string `json:"description" validate:"omitempty,required,min=3,max=100"`
}

func (req CreateBoardRequest) ToInput(userId string) (CreateBoardInput, error) {
	if userIdUUID, err := uuid.Parse(userId); err != nil {
		return CreateBoardInput{}, err
	} else {
		return CreateBoardInput{
			Name:        req.Name,
			Description: req.Description,
			UserId:      userIdUUID,
		}, nil
	}
}

type CreateBoardInput struct {
	Name        *string
	Description *string
	UserId      uuid.UUID
}

// Responses
type BoardResponse struct {
	Id          uuid.UUID `json:"id"`
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
	UserId      uuid.UUID `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GetBoardsResponse struct {
	Boards []BoardResponse `json:"boards"`
}
