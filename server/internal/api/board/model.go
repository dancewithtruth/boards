package board

import (
	"time"

	"github.com/google/uuid"
)

type Board struct {
	Id          uuid.UUID
	Name        *string
	Description *string
	UserId      uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
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

type Boards []Board

func (b Boards) ToDto() GetBoardsResponse {
	boardResponses := make([]BoardResponse, len(b))
	for i, board := range b {
		boardResponses[i] = BoardResponse{
			Id:          board.Id,
			Name:        board.Name,
			Description: board.Description,
			UserId:      board.UserId,
			CreatedAt:   board.CreatedAt,
			UpdatedAt:   board.UpdatedAt,
		}
	}
	return GetBoardsResponse{
		Boards: boardResponses,
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
