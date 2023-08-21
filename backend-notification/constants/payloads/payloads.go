package payloads

type Invite struct {
	ID string `json:"id"`
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type EmailVerification struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Code  string `json:"code"`
}
