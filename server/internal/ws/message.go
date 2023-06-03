package ws

import "encoding/json"

const (
	// Message events
	EventBoardConnectUser = "board.connect_user"
	EventPostCreate       = "post.create"

	// Close reasons
	CloseReasonBadEvent         = "The event is not the correct type. Please ensure event is a string."
	CloseReasonMissingEvent     = "The event is missing. Please include a event."
	CloseReasonBadParams        = "The params have incorrect field types. Please ensure params have the correct types."
	CloseReasonUnsupportedEvent = "The event is not supported. Please make sure you have the right event."

	// Error messages
	ErrorMessageInternalCreatePost = "There was an issue creating a post."
)

// Request is a struct that describes the shape of every message request
type Request struct {
	Event  string          `json:"event"`
	Params json.RawMessage `json:"params"`
}

// Response is a struct that describes the shape of every message response. If a request is handled and
// encounters an error, then display the appropriate error using ErrorMessage.
type Response struct {
	Event        string `json:"event"`
	Success      bool   `json:"sucess"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// ResponseUserConnect is a struct that describes the shape of a connect user response
type ResponseUserConnect struct {
	Response
	Result ResultUserConnect `json:"result"`
}

type ResultUserConnect struct {
	NewUser       string   `json:"new_user"`
	ExistingUsers []string `json:"existing_users"`
}
