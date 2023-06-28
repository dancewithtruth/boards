package post

import "github.com/Wave-95/boards/backend-core/pkg/validator"

// CreatePostInput defines the structure of a request to create a post
type CreatePostInput struct {
	UserID  string `json:"user_id" validate:"required,uuid"`
	BoardID string `json:"board_id" validate:"required,uuid"`
	Content string `json:"content"`
	PosX    int    `json:"pos_x" validate:"required,min=0"`
	PosY    int    `json:"pos_y" validate:"required,min=0"`
	Color   string `json:"color" validate:"required,min=7,max=7"`
	Height  int    `json:"height" validate:"min=0"`
	ZIndex  int    `json:"z_index" validate:"min=1"`
}

// Validate validates the create post input.
func (i *CreatePostInput) Validate() error {
	validator := validator.New()
	return validator.Struct(i)
}

// UpdatePostInput defines the structure of a request to create a post
type UpdatePostInput struct {
	ID      string  `json:"id" validate:"required,uuid"`
	Content *string `json:"content"`
	PosX    *int    `json:"pos_x" validate:"omitempty,min=0"`
	PosY    *int    `json:"pos_y" validate:"omitempty,min=0"`
	Color   *string `json:"color" validate:"omitempty,min=7,max=7"`
	Height  *int    `json:"height" validate:"omitempty,min=0"`
	ZIndex  *int    `json:"z_index" validate:"omitempty,min=1"`
}

// Validate validates the update post payload.
func (i *UpdatePostInput) Validate() error {
	validator := validator.New()
	return validator.Struct(i)
}
