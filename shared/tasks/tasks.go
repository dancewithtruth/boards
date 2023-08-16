package tasks

const (
	EmailInvite = "task_email_invite"
)

type Message struct {
	Task    string `json:"task"`
	Payload any    `json:"payload"`
}
