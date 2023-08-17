package handlers

type CreateEmailVerificationInput struct {
	UserID string `json:"user_id"`
	Code   string `json:"code"`
}
