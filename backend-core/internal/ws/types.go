package ws

import (
	"encoding/json"

	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/Wave-95/boards/backend-core/internal/post"
	"github.com/google/uuid"
)

const (
	// Events

	// EventUserAuthenticate is when a user authenticates.
	EventUserAuthenticate = "user.authenticate"

	// EventBoardConnect is when a board is connected.
	EventBoardConnect = "board.connect"

	// EventBoardDisconnect is when a board is disconnected.
	EventBoardDisconnect = "board.disconnect"

	// EventPostCreate is when a post is created.
	EventPostCreate = "post.create"

	// EventPostUpdate is when a post is updated.
	EventPostUpdate = "post.update"

	// EventPostDelete is when a post is deleted.
	EventPostDelete = "post.delete"

	// EventPostFocus is when a post receives focus.
	EventPostFocus = "post.focus"

	// EventPostGroupUpdate is when a post group is updated.
	EventPostGroupUpdate = "post_group.update"

	// EventPostGroupDelete is when a post group is deleted.
	EventPostGroupDelete = "post_group.delete"

	// Close Reasons

	// CloseReasonMissingEvent indicates that the event field is missing.
	CloseReasonMissingEvent = "The event field is missing."

	// CloseReasonUnsupportedEvent indicates that the event is unsupported.
	CloseReasonUnsupportedEvent = "The event is unsupported."

	// CloseReasonBadEvent indicates that the event field has an incorrect type.
	CloseReasonBadEvent = "The event field is an incorrect type."

	// CloseReasonBadParams indicates that the params have incorrect field types.
	CloseReasonBadParams = "The params have incorrect field types."

	// CloseReasonInternalServer indicates an internal server error.
	CloseReasonInternalServer = "Internal server error."

	// CloseReasonUnauthorized indicates an unauthorized request.
	CloseReasonUnauthorized = "Unauthorized."

	// Error Messages

	// ErrMsgInvalidJwt indicates an invalid JWT token.
	ErrMsgInvalidJwt = "Invalid JWT token supplied."

	// ErrMsgBoardNotFound indicates that a board was not found.
	ErrMsgBoardNotFound = "Board not found."

	// ErrMsgUnauthorized indicates an unauthorized request.
	ErrMsgUnauthorized = "Unauthorized."

	// ErrMsgInternalServer indicates an internal server error.
	ErrMsgInternalServer = "Internal server error."
)

// Request is a struct that describes the shape of every message request.
type Request struct {
	Event  string          `json:"event"`
	Params json.RawMessage `json:"params"`
}

// RequestBoardConnect represents a request to connect to a board.
type RequestBoardConnect struct {
	Event  string             `json:"event"`
	Params ParamsBoardConnect `json:"params"`
}

// ParamsBoardConnect contains the parameters for board connection.
type ParamsBoardConnect struct {
	BoardID string `json:"board_id"`
}

// RequestUserAuthenticate represents a request to authenticate a user.
type RequestUserAuthenticate struct {
	Event  string                 `json:"event"`
	Params ParamsUserAuthenticate `json:"params"`
}

// ParamsUserAuthenticate contains the parameters for user authentication.
type ParamsUserAuthenticate struct {
	Jwt string `json:"jwt"`
}

// RequestPostCreate represents a request to create a post.
type RequestPostCreate struct {
	Event  string           `json:"event"`
	Params ParamsPostCreate `json:"params"`
}

// ParamsPostCreate contains the parameters for post creation.
type ParamsPostCreate struct {
	BoardID     string  `json:"board_id" validate:"required,uuid"`
	Content     string  `json:"content"`
	PosX        int     `json:"pos_x" validate:"required,min=0"`
	PosY        int     `json:"pos_y" validate:"required,min=0"`
	Color       string  `json:"color" validate:"required,min=7,max=7"`
	Height      int     `json:"height" validate:"min=0"`
	ZIndex      int     `json:"z_index" validate:"min=1"`
	PostOrder   float64 `json:"post_order"`
	PostGroupID string  `json:"post_group_id"`
}

// RequestPostUpdate represents a request to update a post.
type RequestPostUpdate struct {
	Event  string           `json:"event"`
	Params ParamsPostUpdate `json:"params"`
}

// ParamsPostUpdate contains the parameters for post update.
type ParamsPostUpdate struct {
	BoardID string `json:"board_id" validate:"required,uuid"`
	post.UpdatePostInput
}

// RequestPostDelete represents a request to delete a post.
type RequestPostDelete struct {
	Event  string           `json:"event"`
	Params ParamsPostDelete `json:"params"`
}

// ParamsPostDelete contains the parameters for post deletion.
type ParamsPostDelete struct {
	PostID  string `json:"post_id" validate:"required,uuid"`
	BoardID string `json:"board_id" validate:"required,uuid"`
}

// RequestPostFocus represents a request to focus on a post.
type RequestPostFocus struct {
	Event  string           `json:"event"`
	Params ParamsPostDelete `json:"params"`
}

// ParamsPostFocus contains the parameters for post focus.
type ParamsPostFocus struct {
	ID      string `json:"id" validate:"required,uuid"`
	BoardID string `json:"board_id" validate:"required,uuid"`
}

// RequestPostGroupUpdate represents a request to update a post group.
type RequestPostGroupUpdate struct {
	Event  string                `json:"event"`
	Params ParamsPostGroupUpdate `json:"params"`
}

// ParamsPostGroupUpdate contains the parameters for updating a post group.
type ParamsPostGroupUpdate struct {
	BoardID string `json:"board_id" validate:"required,uuid"`
	post.UpdatePostGroupInput
}

// RequestPostGroupDelete represents a request to delete a post group.
type RequestPostGroupDelete struct {
	Event  string                `json:"event"`
	Params ParamsPostGroupDelete `json:"params"`
}

// ParamsPostGroupDelete contains the parameters for deleting a post group.
type ParamsPostGroupDelete struct {
	PostGroupID string `json:"post_group_id" validate:"required,uuid"`
}

// ResponseBase represents the base response structure.
type ResponseBase struct {
	Event        string `json:"event"`
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// ResponseUserAuthenticate represents the response for user authentication.
type ResponseUserAuthenticate struct {
	Event        string                 `json:"event"`
	Success      bool                   `json:"success"`
	Result       ResultUserAuthenticate `json:"result,omitempty"`
	ErrorMessage string                 `json:"error_message,omitempty"`
}

// ResultUserAuthenticate contains the result of user authentication.
type ResultUserAuthenticate struct {
	User models.User `json:"user"`
}

// ResponseBoardConnect represents the response for board connection.
type ResponseBoardConnect struct {
	ResponseBase
	Result ResultBoardConnect `json:"result,omitempty"`
}

// ResultBoardConnect contains the result of board connection.
type ResultBoardConnect struct {
	BoardID        string        `json:"board_id"`
	NewUser        models.User   `json:"new_user"`
	ConnectedUsers []models.User `json:"connected_users"`
}

// ResponsePostCreate represents the response for creating a new post.
type ResponsePostCreate struct {
	ResponseBase
	Result ResultPostCreate `json:"result,omitempty"`
}

type ResultPostCreate struct {
	Post      models.Post      `json:"post,omitempty"`
	PostGroup models.PostGroup `json:"post_group,omitempty"`
}

// ResponsePostUpdate represents the response for updating a new post.
type ResponsePostUpdate struct {
	ResponseBase
	Result models.Post `json:"result,omitempty"`
}

// ResponsePostDelete represents the response for post deletion.
type ResponsePostDelete struct {
	ResponseBase
	Result models.Post `json:"result,omitempty"`
}

// ResponsePostFocus represents the response for post focusing.
type ResponsePostFocus struct {
	ResponseBase
	Result ResultPostFocus `json:"result,omitempty"`
}

// ResultPostFocus contains the result of post focusing.
type ResultPostFocus struct {
	Post models.Post `json:"post"`
	User models.User `json:"user"`
}

// ResponsePostGroup represents the response for post group.
type ResponsePostGroup struct {
	ResponseBase
	Result models.PostGroup `json:"result,omitempty"`
}

// ResponsePostGroupDeleted represents the response for post group.
type ResponsePostGroupDeleted struct {
	ResponseBase
	Result struct {
		ID uuid.UUID `json:"id"`
	} `json:"result,omitempty"`
}

// ResponseUserDisconnect represents the response for user disconnection.
type ResponseUserDisconnect struct {
	ResponseBase
	Result ResultUserDisconnect `json:"result,omitempty"`
}

// ResultUserDisconnect contains the result of user disconnection.
type ResultUserDisconnect struct {
	UserID string `json:"user_id"`
}
