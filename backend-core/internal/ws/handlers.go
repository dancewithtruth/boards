package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Wave-95/boards/backend-core/internal/board"
	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/Wave-95/boards/backend-core/internal/post"
	"github.com/Wave-95/boards/backend-core/pkg/logger"
	"github.com/Wave-95/boards/backend-core/pkg/validator"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Allow all origins
			return true
		},
	}
)

func (ws *WebSocket) HandleConnection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logger.FromContext(ctx)

	// Upgrade connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Errorf("handler: failed to upgrade connection: %v", err)
		logger.Info("request:", r)
		return
	}

	client := Client{
		boards:        make(map[string]Board),
		subscriptions: make(map[string]chan bool),
		conn:          conn,
		send:          make(chan []byte, 256),
		ws:            ws,
	}

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

// handleUserAuthenticate will authenticate the user and store a userID in the Client.
func handleUserAuthenticate(c *Client, msgReq Request) {
	// Unmarshal params
	var params ParamsUserAuthenticate
	if err := unmarshalParams(msgReq, &params, c); err != nil {
		return
	}
	// Verify user
	userID, _ := c.ws.jwtService.VerifyToken(params.Jwt)
	user, err := c.ws.userService.GetUser(context.Background(), userID)
	user.Password = nil
	msgRes := ResponseUserAuthenticate{Event: EventUserAuthenticate}
	if err != nil {
		// Prepare error response
		msgRes.Success = false
		msgRes.ErrorMessage = ErrMsgInvalidJwt
	} else {
		// Prepare success response
		msgRes.Success = true
		msgRes.Result = ResultUserAuthenticate{User: user}
		// Assign userID to client struct
		c.user = &user
	}

	// Send response
	msgResBytes, err := json.Marshal(msgRes)
	if err := handleMarshalError(err, "handleUserAuthenticate", c); err != nil {
		return
	}
	c.send <- msgResBytes
}

// handleBoardConnect will subscribe a client to a board channel and, if successful, will broadcast
// the connection to all subscribers
func handleBoardConnect(c *Client, msgReq Request) {
	// Check if user is authenticated
	user := c.user
	if user == nil {
		closeConnection(c, websocket.ClosePolicyViolation, CloseReasonUnauthorized)
		return
	}
	// Unmarshal params
	var params ParamsBoardConnect
	if err := unmarshalParams(msgReq, &params, c); err != nil {
		return
	}
	// Check if user has access to board
	boardID := params.BoardID
	boardWithMembers, err := c.ws.boardService.GetBoardWithMembers(context.Background(), boardID)
	var msgRes ResponseBoardConnect
	// If no access, return error response
	if err != nil || !board.UserHasAccess(boardWithMembers, user.ID.String()) {
		msgRes = ResponseBoardConnect{
			ResponseBase: ResponseBase{
				Event:        EventBoardConnect,
				Success:      false,
				ErrorMessage: ErrMsgBoardNotFound},
		}
	} else {
		rdb := c.ws.rdb

		go c.subscribe(boardID)

		if err := setUser(rdb, boardID, *user); err != nil {
			fmt.Printf("Issue setting user using HSet: %v", err)
			closeConnection(c, websocket.CloseProtocolError, CloseReasonInternalServer)
		}

		mp, err := getUsers(rdb, boardID)
		if err != nil {
			fmt.Printf("Issue getting user using HGetAll: %v", err)
			closeConnection(c, websocket.CloseProtocolError, CloseReasonInternalServer)
		}

		connectedUsers, err := formatConnectedUsers(mp)
		if err != nil {
			fmt.Printf("Issue formatting connected users: %v", err)
			closeConnection(c, websocket.CloseProtocolError, CloseReasonInternalServer)
		}

		// Broacast successful message response
		msgRes = ResponseBoardConnect{
			ResponseBase: ResponseBase{
				Event:   EventBoardConnect,
				Success: true,
			},
			Result: ResultBoardConnect{
				BoardID:        boardID,
				NewUser:        *user,
				ConnectedUsers: connectedUsers,
			},
		}
	}
	msgResBytes, err := json.Marshal(msgRes)
	if err := handleMarshalError(err, "handleBoardConnect", c); err != nil {
		return
	}
	if msgRes.Success {
		c.ws.rdb.Publish(context.Background(), boardID, msgResBytes)
	}
	c.send <- msgResBytes
}

func handlePostCreate(c *Client, msgReq Request) {
	// Authenticate user
	user := c.user
	if user.ID.String() == "" {
		closeConnection(c, websocket.ClosePolicyViolation, CloseReasonUnauthorized)
		return
	}
	// Unmarshal message request
	var params ParamsPostCreate
	if err := unmarshalParams(msgReq, &params, c); err != nil {
		return
	}
	// Prepare create post input
	createPostInput := post.CreatePostInput{
		UserID:      user.ID.String(),
		BoardID:     params.BoardID,
		Content:     params.Content,
		PosX:        params.PosX,
		PosY:        params.PosY,
		Color:       params.Color,
		Height:      params.Height,
		ZIndex:      params.ZIndex,
		PostOrder:   params.PostOrder,
		PostGroupID: params.PostGroupID,
	}
	post, err := c.ws.postService.CreatePost(context.Background(), createPostInput)
	if err != nil {
		switch {
		case validator.IsValidationError(err):
			validationErrMsg := validator.GetValidationErrMsg(createPostInput, err)
			sendErrorMessage(c, buildErrorResponse(msgReq, validationErrMsg))
		default:
			sendErrorMessage(c, buildErrorResponse(msgReq, ErrMsgInternalServer))
		}
		return
	}
	// Get parent post group
	postGroup, err := c.ws.postService.GetPostGroup(context.Background(), post.PostGroupID.String())
	if err != nil {
		log.Printf("handler: failed to get post group: %v", err)
		sendErrorMessage(c, buildErrorResponse(msgReq, ErrMsgInternalServer))
		return
	}
	// Prepare message response
	msgRes := ResponsePostCreate{
		ResponseBase: ResponseBase{
			Event:   msgReq.Event,
			Success: true,
		},
		Result: ResultPostCreate{
			Post:      post,
			PostGroup: postGroup,
		},
	}
	msgResBytes, err := json.Marshal(msgRes)
	if err := handleMarshalError(err, "handlePostCreate", c); err != nil {
		return
	}
	// Broadcast message response
	c.ws.rdb.Publish(context.Background(), params.BoardID, msgResBytes)
}

func handlePostFocus(c *Client, msgReq Request) {
	user := c.user
	if user.ID.String() == "" {
		closeConnection(c, websocket.ClosePolicyViolation, CloseReasonUnauthorized)
		return
	}
	var params ParamsPostFocus
	if err := unmarshalParams(msgReq, &params, c); err != nil {
		return
	}
	postID := params.ID
	post, err := c.ws.postService.GetPost(context.Background(), postID)
	if err != nil {
		log.Printf("handler: failed to get post during post focus request: %v", err)
		sendErrorMessage(c, buildErrorResponse(msgReq, ErrMsgInternalServer))
	}
	postGroup, err := c.ws.postService.GetPostGroup(context.Background(), post.PostGroupID.String())
	if err != nil {
		log.Printf("handler: failed to get post group during post focus request: %v", err)
		sendErrorMessage(c, buildErrorResponse(msgReq, ErrMsgInternalServer))
	}
	boardID := postGroup.BoardID.String()
	msgRes := ResponsePostFocus{
		ResponseBase: ResponseBase{
			Event:   msgReq.Event,
			Success: true,
		},
		Result: ResultPostFocus{
			Post: post,
			User: *c.user,
		},
	}
	msgResBytes, err := json.Marshal(msgRes)

	if err := handleMarshalError(err, "handlePostFocus", c); err != nil {
		return
	}
	c.ws.rdb.Publish(context.Background(), boardID, msgResBytes)
}

func handlePostUpdate(c *Client, msgReq Request) {
	// Authenticate user
	user := c.user
	if user.ID.String() == "" {
		closeConnection(c, websocket.ClosePolicyViolation, CloseReasonUnauthorized)
		return
	}
	// Unmarshal message request
	var params ParamsPostUpdate
	if err := unmarshalParams(msgReq, &params, c); err != nil {
		return
	}
	existingPost, err := c.ws.postService.GetPost(context.Background(), params.UpdatePostInput.ID)
	if err != nil {
		sendErrorMessage(c, buildErrorResponse(msgReq, ErrMsgInternalServer))
	}
	postGroup, err := c.ws.postService.GetPostGroup(context.Background(), existingPost.PostGroupID.String())
	if err != nil {
		sendErrorMessage(c, buildErrorResponse(msgReq, ErrMsgInternalServer))
	}
	boardID := postGroup.BoardID.String()
	// Prepare update post input
	updatePostInput := post.UpdatePostInput{
		ID:          params.ID,
		Content:     params.Content,
		Color:       params.Color,
		Height:      params.Height,
		PostOrder:   params.PostOrder,
		PostGroupID: params.PostGroupID,
	}
	updatedPost, err := c.ws.postService.UpdatePost(context.Background(), updatePostInput)
	if err != nil {
		switch {
		case validator.IsValidationError(err):
			validationErrMsg := validator.GetValidationErrMsg(updatePostInput, err)
			sendErrorMessage(c, buildErrorResponse(msgReq, validationErrMsg))
		default:
			sendErrorMessage(c, buildErrorResponse(msgReq, ErrMsgInternalServer))
		}
		return
	}
	// Prepare update post message response
	msgRes := ResponsePostUpdate{
		ResponseBase: ResponseBase{
			Event:   msgReq.Event,
			Success: true,
		},
		Result: ResultPostUpdate{
			OldPost:     existingPost,
			UpdatedPost: updatedPost,
		},
	}
	msgResBytes, err := json.Marshal(msgRes)
	if err := handleMarshalError(err, "handlePostUpdate", c); err != nil {
		return
	}
	// Broadcast message response
	c.ws.rdb.Publish(context.Background(), boardID, msgResBytes)
}

// handlePostDetach detaches a post from its original post group and creates then assigns it to
// the newly created post group.
func handlePostDetach(c *Client, msgReq Request) {
	// Authenticate user
	user := c.user
	if user.ID.String() == "" {
		closeConnection(c, websocket.ClosePolicyViolation, CloseReasonUnauthorized)
		return
	}
	// Unmarshal message request
	var params ParamsPostDetach
	if err := unmarshalParams(msgReq, &params, c); err != nil {
		return
	}

	postID := params.ID
	existingPost, err := c.ws.postService.GetPost(context.Background(), postID)
	if err != nil {
		log.Printf("handler: failed to get post during post detach request: %v", err)
		sendErrorMessage(c, buildErrorResponse(msgReq, ErrMsgInternalServer))
	}
	existingPostGroup, err := c.ws.postService.GetPostGroup(context.Background(), existingPost.PostGroupID.String())
	if err != nil {
		log.Printf("handler: failed to get post group during post detach request: %v", err)
		sendErrorMessage(c, buildErrorResponse(msgReq, ErrMsgInternalServer))
	}
	boardID := existingPostGroup.BoardID.String()
	// Create new post group
	createPostGroupInput := post.CreatePostGroupInput{
		BoardID: boardID,
		PosX:    params.PosX,
		PosY:    params.PosY,
		ZIndex:  params.ZIndex,
	}
	newPostGroup, err := c.ws.postService.CreatePostGroup(context.Background(), createPostGroupInput)
	if err != nil {
		log.Printf("handler: failed to create a new post group during post detach request: %v", err)
		sendErrorMessage(c, buildErrorResponse(msgReq, ErrMsgInternalServer))
	}
	newPostGroupID := newPostGroup.ID.String()
	updatePostInput := post.UpdatePostInput{
		ID:          existingPost.ID.String(),
		PostGroupID: &newPostGroupID,
	}
	updatedPost, err := c.ws.postService.UpdatePost(context.Background(), updatePostInput)
	if err != nil {
		switch {
		case validator.IsValidationError(err):
			validationErrMsg := validator.GetValidationErrMsg(updatePostInput, err)
			sendErrorMessage(c, buildErrorResponse(msgReq, validationErrMsg))
		default:
			sendErrorMessage(c, buildErrorResponse(msgReq, ErrMsgInternalServer))
		}
		return
	}
	// Prepare message response
	msgRes := ResponsePostDetach{
		ResponseBase: ResponseBase{
			Event:   msgReq.Event,
			Success: true,
		},
		Result: ResultPostDetach{
			OldPost:     existingPost,
			UpdatedPost: updatedPost,
			PostGroup:   newPostGroup,
		},
	}
	msgResBytes, err := json.Marshal(msgRes)
	if err := handleMarshalError(err, "handlePostDetach", c); err != nil {
		return
	}
	// Broadcast message response
	c.ws.rdb.Publish(context.Background(), boardID, msgResBytes)
}

func handlePostDelete(c *Client, msgReq Request) {
	user := c.user
	if user.ID.String() == "" {
		closeConnection(c, websocket.ClosePolicyViolation, CloseReasonUnauthorized)
		return
	}
	var params ParamsPostDelete
	if err := unmarshalParams(msgReq, &params, c); err != nil {
		return
	}
	postID := params.PostID
	post, err := c.ws.postService.GetPost(context.Background(), postID)
	if err != nil {
		log.Printf("handler: failed to get post during post delete request: %v", err)
		sendErrorMessage(c, buildErrorResponse(msgReq, ErrMsgInternalServer))
	}
	postGroup, err := c.ws.postService.GetPostGroup(context.Background(), post.PostGroupID.String())
	if err != nil {
		log.Printf("handler: failed to get post group during post delete request: %v", err)
		sendErrorMessage(c, buildErrorResponse(msgReq, ErrMsgInternalServer))
	}
	boardID := postGroup.BoardID.String()
	if err := c.ws.postService.DeletePost(context.Background(), postID); err != nil {
		sendErrorMessage(c, buildErrorResponse(msgReq, ErrMsgInternalServer))
		return
	}
	msgRes := ResponsePostDelete{
		ResponseBase: ResponseBase{
			Event:   msgReq.Event,
			Success: true,
		},
		Result: post,
	}
	msgResBytes, err := json.Marshal(msgRes)
	if err := handleMarshalError(err, "handlePostDelete", c); err != nil {
		return
	}
	c.ws.rdb.Publish(context.Background(), boardID, msgResBytes)
}

// handlePostGroupUpdate handles a message request to update a post group.
func handlePostGroupUpdate(c *Client, msgReq Request) {
	user := c.user
	if user.ID.String() == "" {
		closeConnection(c, websocket.ClosePolicyViolation, CloseReasonUnauthorized)
		return
	}
	var params ParamsPostGroupUpdate
	if err := unmarshalParams(msgReq, &params, c); err != nil {
		return
	}
	boardID := params.BoardID
	updatePostInput := post.UpdatePostGroupInput{
		ID:     params.ID,
		Title:  params.Title,
		PosX:   params.PosX,
		PosY:   params.PosY,
		ZIndex: params.ZIndex,
	}
	postGroup, err := c.ws.postService.UpdatePostGroup(context.Background(), updatePostInput)
	if err != nil {
		switch {
		case validator.IsValidationError(err):
			validationErrMsg := validator.GetValidationErrMsg(updatePostInput, err)
			sendErrorMessage(c, buildErrorResponse(msgReq, validationErrMsg))
		default:
			sendErrorMessage(c, buildErrorResponse(msgReq, ErrMsgInternalServer))
		}
		return
	}
	msgRes := ResponsePostGroup{
		ResponseBase: ResponseBase{
			Event:   msgReq.Event,
			Success: true,
		},
		Result: postGroup,
	}
	msgResBytes, err := json.Marshal(msgRes)
	if err := handleMarshalError(err, "handlePostUpdate", c); err != nil {
		return
	}
	c.ws.rdb.Publish(context.Background(), boardID, msgResBytes)

}

// handlePostGroupDelete handles a message request to delete a post group.
func handlePostGroupDelete(c *Client, msgReq Request) {
	// Authenticate user
	user := c.user
	if user.ID.String() == "" {
		closeConnection(c, websocket.ClosePolicyViolation, CloseReasonUnauthorized)
		return
	}
	// Unmarshal request
	var params ParamsPostGroupDelete
	if err := unmarshalParams(msgReq, &params, c); err != nil {
		return
	}
	// Get post group and check if user has write permissions
	postGroupID := params.PostGroupID
	postGroup, err := c.ws.postService.GetPostGroup(context.Background(), postGroupID)
	if err != nil {
		log.Printf("handler: failed to get post group: %v", err)
		sendErrorMessage(c, buildErrorResponse(msgReq, ErrMsgInternalServer))
		return
	}
	boardID := postGroup.BoardID.String()
	// Delete post group
	err = c.ws.postService.DeletePostGroup(context.Background(), postGroupID)
	if err != nil {
		log.Printf("handler: failed to delete post group: %v", err)
		sendErrorMessage(c, buildErrorResponse(msgReq, ErrMsgInternalServer))
		return
	}
	// Prepare response
	msgRes := ResponsePostGroupDeleted{
		ResponseBase: ResponseBase{
			Event:   msgReq.Event,
			Success: true,
		},
		Result: struct {
			ID uuid.UUID `json:"id"`
		}{postGroup.ID},
	}
	msgResBytes, err := json.Marshal(msgRes)
	if err := handleMarshalError(err, "handlePostUpdate", c); err != nil {
		return
	}
	// Broadcast response
	c.ws.rdb.Publish(context.Background(), boardID, msgResBytes)
}

// unmarshalParams is a helper function that unmarshals a message request's params and sends
// out a close connection message if any errors are encountered.
func unmarshalParams(msgReq Request, v any, c *Client) error {
	err := json.Unmarshal(msgReq.Params, v)
	if err != nil {
		closeConnection(c, websocket.CloseInvalidFramePayloadData, CloseReasonBadParams)
		return err
	}
	return nil
}

// handleMarshalError checks to see if there are any errors when marshalling the WebSocket message response into JSON.
// If there is an issue, it will close the connection with an internal server error close reason.
func handleMarshalError(err error, handlerName string, c *Client) error {
	if err != nil {
		log.Printf("%s: Failed to marshal response into JSON: %v", handlerName, err)
		closeConnection(c, websocket.CloseProtocolError, CloseReasonInternalServer)
		return err
	}
	return nil
}

// closeConnection is a helper function to write a control message with a status code and close reason text.
func closeConnection(c *Client, statusCode int, text string) {
	err := c.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(statusCode, text), time.Now().Add(writeWait))
	if err != nil {
		log.Printf("Failed to write control message: %v", err)
	}
}

// formatConnectedUsers formats the map of connected users into an array of users
func formatConnectedUsers(mp map[string]string) ([]models.User, error) {
	users := make([]models.User, len(mp))
	i := 0
	for _, val := range mp {
		var user models.User
		err := json.Unmarshal([]byte(val), &user)
		if err != nil {
			return []models.User{}, err
		}
		users[i] = user
		i++
	}
	return users, nil
}
