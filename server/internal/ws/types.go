package ws

import (
	"encoding/json"

	"github.com/Wave-95/boards/server/internal/models"
	"github.com/Wave-95/boards/server/internal/post"
)

const (
	// Events
	EventUserAuthenticate = "user.authenticate"
	EventBoardConnect     = "board.connect"
	EventPostCreate       = "post.create"
	EventPostUpdate       = "post.update"
	EventPostDelete       = "post.delete"
	EventPostFocus        = "post.focus"

	// Close reasons
	CloseReasonMissingEvent     = "The event field is missing."
	CloseReasonUnsupportedEvent = "The event is unsupported."
	CloseReasonBadEvent         = "The event field is an incorrect type."
	CloseReasonBadParams        = "The params have incorrect field types. Please ensure params have the correct types."
	CloseReasonInternalServer   = "Internal server error."
	CloseReasonUnauthorized     = "Unauthorized."

	// Error messages
	ErrMsgInvalidJwt     = "Invalid JWT token supplied."
	ErrMsgBoardNotFound  = "Board not found."
	ErrMsgUnauthorized   = "Unauthorized."
	ErrMsgInternalServer = "Internal server error."
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

type RequestPostCreate struct {
	Event  string           `json:"event"`
	Params ParamsPostCreate `json:"params"`
}

type ParamsPostCreate struct {
	BoardId string `json:"board_id" validate:"required,uuid"`
	Content string `json:"content"`
	PosX    int    `json:"pos_x" validate:"required,min=0"`
	PosY    int    `json:"pos_y" validate:"required,min=0"`
	Color   string `json:"color" validate:"required,min=7,max=7"`
	Height  int    `json:"height" validate:"min=0"`
	ZIndex  int    `json:"z_index" validate:"min=1"`
}

type RequestPostUpdate struct {
	Event  string           `json:"event"`
	Params ParamsPostUpdate `json:"params"`
}

type ParamsPostUpdate struct {
	BoardId string `json:"board_id" validate:"required,uuid"`
	post.UpdatePostInput
}

type RequestPostDelete struct {
	Event  string           `json:"event"`
	Params ParamsPostDelete `json:"params"`
}

type ParamsPostDelete struct {
	PostId  string `json:"post_id" validate:"required,uuid"`
	BoardId string `json:"board_id" validate:"required,uuid"`
}

type RequestPostFocus struct {
	Event  string           `json:"event"`
	Params ParamsPostDelete `json:"params"`
}

type ParamsPostFocus struct {
	Id      string `json:"id" validate:"required,uuid"`
	BoardId string `json:"board_id" validate:"required,uuid"`
}

// Responses

type ResponseBase struct {
	Event        string `json:"event"`
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type ResponseUserAuthenticate struct {
	Event        string                 `json:"event"`
	Success      bool                   `json:"success"`
	Result       ResultUserAuthenticate `json:"result,omitempty"`
	ErrorMessage string                 `json:"error_message,omitempty"`
}
type ResultUserAuthenticate struct {
	User models.User `json:"user"`
}
type ResponseBoardConnect struct {
	ResponseBase
	Result ResultBoardConnect `json:"result,omitempty"`
}
type ResultBoardConnect struct {
	BoardId        string        `json:"board_id"`
	NewUser        models.User   `json:"new_user"`
	ConnectedUsers []models.User `json:"connected_users"`
}

type ResponsePost struct {
	ResponseBase
	Result models.Post `json:"result,omitempty"`
}

type ResponsePostDelete struct {
	ResponseBase
	Result ResultPostDelete `json:"result,omitempty"`
}

type ResultPostDelete struct {
	PostId  string `json:"post_id"`
	BoardId string `json:"board_id"`
}

type ResponsePostFocus struct {
	ResponseBase
	Result ResultPostFocus `json:"result,omitempty"`
}

type ResultPostFocus struct {
	Id      string      `json:"id"`
	BoardId string      `json:"board_id"`
	User    models.User `json:"user"`
}
