package ws

import "encoding/json"

const (
	// Message types
	TypeConnectUser         = "CONNECT_USER"
	TypeCreatePost          = "CREATE_POST"
	TypeBadRequest          = "BAD_REQUEST"
	TypeInvalidRequest      = "INVALID_REQUEST"
	TypeInternalServerError = "INTERNAL_SERVER_ERROR"

	// Error messages
	ErrorMessageBadRequest         = "Bad message request. Please make sure your field types are correct."
	ErrorMessageBadType            = "Message request type not supported. Please make sure you have the right type."
	ErrorMessageInternalCreatePost = "There was an issue creating a post."
)

// MessageRequest is a struct that describes the shape of every message request
type MessageRequest struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// MessageResponse is a struct that describes the shape of every message response. If a request is handled and
// encounters an error, then display the appropriate error using ErrorMessage.
type MessageResponse struct {
	Type         string `json:"type"`
	ErrorMessage string `json:"error,omitempty"`
	Sender       string `json:"sender,omitempty"`
}

// MessageResponseConnectUser is a struct that describes the shape of a connect user response
type MessageResponseConnectUser struct {
	MessageResponse
	Data DataConnectUser `json:"data"`
}

type DataConnectUser struct {
	NewUser       string   `json:"new_user"`
	ExistingUsers []string `json:"existing_users"`
}
