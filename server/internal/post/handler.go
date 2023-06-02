package post

// Inputs
type CreatePostInput struct {
	UserId  string `json:"user_id" validate:"required,uuid"`
	BoardId string `json:"board_id" validate:"required,uuid"`
	Color   string `json:"color" validate:"required,min=7,max=7"`
	Top     int    `json:"top" validate:"required,min=0"`
	Left    int    `json:"left" validate:"required,min=0"`
}
