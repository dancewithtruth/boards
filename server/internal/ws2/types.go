package ws2

import "encoding/json"

const (
	// Events
	EventConnectUser = "board.connect_user"
)

// Requests

// Request is a struct that describes the shape of every message request
type Request struct {
	Event  string          `json:"event"`
	Params json.RawMessage `json:"params"`
}

type RequestConnectUser struct {
	Event  string            `json:"event"`
	Params ParamsConnectUser `json:"params"`
}

type ParamsConnectUser struct {
	Jwt     string `json:"jwt"`
	BoardId string `json:"board_id"`
}

// Responses

type ResponseConnectUser struct {
	Event   string            `json:"event"`
	Success bool              `json:"success"`
	Result  ResultConnectUser `json:"result"`
}

type ResultConnectUser struct {
	NewUser       string   `json:"new_user"`
	ExistingUsers []string `json:"existing_users"`
	BoardId       string   `json:"board_id"`
}
