package auth

// LoginInput represents the input structure for a login request
type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginDTO represents the response structure for a successful login request
type LoginDTO struct {
	Token string `json:"token"`
}
