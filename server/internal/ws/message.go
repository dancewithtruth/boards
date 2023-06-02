package ws

import "encoding/json"

const (
	// Message types
	TypeConnectUser = "CONNECT_USER"
	TypeCreatePost  = "CREATE_POST"

	// Error messages
	ErrorMessageBadType            = "Bad message type. Please make sure the type field is a string."
	ErrorMessageMissingType        = "Message request type is missing. Please make sure you have specified a type."
	ErrorMessageUnsupportedType    = "Message request type is not supported. Please make sure you have the right type."
	ErrorMessageBadInput           = "Bad payload. Please make sure all the fields are the correct type."
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
	Type         string `json:"type,omitempty"`
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
