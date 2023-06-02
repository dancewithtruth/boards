package auth

// Inputs

type LoginInput struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// DTOs

type LoginDTO struct {
	Token string `json:"token"`
}
