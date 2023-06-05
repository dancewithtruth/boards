package ws

import "encoding/json"

const (
	// Events
	EventUserAuthenticate = "user.authenticate"
	EventBoardConnect     = "board.connect"

	// Close reasons
	CloseReasonMissingEvent     = "The event field is missing."
	CloseReasonUnsupportedEvent = "The event is unsupported."
	CloseReasonBadEvent         = "The event field is an incorrect type."
	CloseReasonBadParams        = "The params have incorrect field types. Please ensure params have the correct types."
	CloseReasonInternalServer   = "Internal server error."
	CloseReasonUnauthorized     = "Unauthorized."

	// Error messages
	ErrMsgInvalidJwt    = "Invalid JWT token supplied."
	ErrMsgBoardNotFound = "Board not found."
)

// Requests

// Request is a struct that describes the shape of every message request
type Request struct {
	Event  string          `json:"event"`
	Params json.RawMessage `json:"params"`
}

type RequestBoardConnect struct {
	Event  string             `json:"event"`
	Params ParamsBoardConnect `json:"params"`
}

type ParamsBoardConnect struct {
	BoardId string `json:"board_id"`
}

type RequestUserAuthenticate struct {
	Event  string                 `json:"event"`
	Params ParamsUserAuthenticate `json:"params"`
}

type ParamsUserAuthenticate struct {
	Jwt string `json:"jwt"`
}

// Responses

type ResponseUserAuthenticate struct {
	Event        string                 `json:"event"`
	Success      bool                   `json:"success"`
	Result       ResultUserAuthenticate `json:"result,omitempty"`
	ErrorMessage string                 `json:"error_message,omitempty"`
}
type ResultUserAuthenticate struct {
	UserId string `json:"user_id"`
}
type ResponseBoardConnect struct {
	Event        string             `json:"event"`
	Success      bool               `json:"success"`
	Result       ResultBoardConnect `json:"result,omitempty"`
	ErrorMessage string             `json:"error_message,omitempty"`
}
type ResultBoardConnect struct {
	BoardId       string   `json:"board_id"`
	UserId        string   `json:"user_id"`
	ExistingUsers []string `json:"existing_users"`
}
