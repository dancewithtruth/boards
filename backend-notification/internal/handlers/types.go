package handlers

type CreateEmailVerificationInput struct {
	UserID string `json:"user_id"`
	Code   string `json:"code"`
}

type InviteResponse struct {
	Sender   User `json:"sender"`
	Receiver User `json:"receiver"`
}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
