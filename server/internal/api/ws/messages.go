package ws

const (
	MessageTypeConnectUser = "USER_CONNECTED"
)

type MessageResponse struct {
	Type   string `json:"type"`
	Error  string `json:"error,omitempty"`
	Sender string `json:"sender"`
}

type MessageResponseConnectUser struct {
	MessageResponse
	Data DataConnectUser `json:"data"`
}

type DataConnectUser struct {
	NewUser       string   `json:"new_user"`
	ExistingUsers []string `json:"existing_users"`
}
