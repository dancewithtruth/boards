package ws

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Wave-95/boards/backend-core/internal/board"
	"github.com/Wave-95/boards/backend-core/internal/post"
	"github.com/Wave-95/boards/backend-core/pkg/logger"
	v "github.com/Wave-95/boards/backend-core/pkg/validator"
	"github.com/go-playground/validator/v10"
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
		return
	}

	client := Client{
		boards: make(map[string]Board),
		conn:   conn,
		send:   make(chan []byte, 256),
		ws:     ws,
	}

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

func handleDestroy(destroy chan string, boardHubs map[string]*Hub) {
	for {
		boardID := <-destroy
		delete(boardHubs, boardID)
	}
}

// handleUserAuthenticate will authenticate the user and store a userID in the Client.
func handleUserAuthenticate(c *Client, msgReq Request) {
	// Unmarshal params
	var params ParamsUserAuthenticate
	err := json.Unmarshal(msgReq.Params, &params)
	if err != nil {
		closeConnection(c, websocket.CloseInvalidFramePayloadData, CloseReasonBadParams)
		return
	}
	// Verify user
	userID, err := c.ws.jwtService.VerifyToken(params.Jwt)
	user, err := c.ws.userService.GetUser(context.Background(), userID)
	user.Password = nil
	var msgRes ResponseUserAuthenticate
	if err != nil {
		// Prepare error response
		msgRes = ResponseUserAuthenticate{
			Event:        EventUserAuthenticate,
			Success:      false,
			ErrorMessage: ErrMsgInvalidJwt,
		}
	} else {
		// Prepare success response
		msgRes = ResponseUserAuthenticate{
			Event:   EventUserAuthenticate,
			Success: true,
			Result: ResultUserAuthenticate{
				User: user,
			},
		}
		// Assign userID to client struct
		c.user = &user
	}

	// Send response
	msgResBytes, err := json.Marshal(msgRes)
	if err != nil {
		log.Printf("handleUserAuthenticate: Failed to marshal response into JSON: %v", err)
		closeConnection(c, websocket.CloseProtocolError, CloseReasonInternalServer)
		return
	}
	c.send <- msgResBytes
}

// handleBoardConnect will attempt to connect a user to the board by registering the client
// to a board hub. If a board hub does not exist, it will be created in this handler. A successful
// board connect event will be broadcasted to all connected clients. The response contains the new
// user that's been connected as well as existing users.
func handleBoardConnect(c *Client, msgReq Request) {
	// Check if user is authenticated
	user := c.user
	if user == nil {
		closeConnection(c, websocket.ClosePolicyViolation, CloseReasonUnauthorized)
		return
	}
	// Unmarshal params
	var params ParamsBoardConnect
	err := json.Unmarshal(msgReq.Params, &params)
	if err != nil {
		closeConnection(c, websocket.CloseInvalidFramePayloadData, CloseReasonBadParams)
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
		// If board does not exist as a hub, create one and run it
		if _, ok := c.ws.boardHubs[boardID]; !ok {
			c.ws.boardHubs[boardID] = newHub(boardID, c.ws.destroy)
			go c.ws.boardHubs[boardID].run()
		}
		// Store write permission on the client's boards map
		boardHub := c.ws.boardHubs[boardID]
		c.boards[boardID] = Board{
			canWrite: true,
		}
		existingUsers := boardHub.listConnectedUsers()
		boardHub.register <- c

		// Broacast successful message response
		msgRes = ResponseBoardConnect{
			ResponseBase: ResponseBase{
				Event:   EventBoardConnect,
				Success: true,
			},
			Result: ResultBoardConnect{
				BoardID:        boardID,
				NewUser:        *user,
				ConnectedUsers: existingUsers,
			},
		}
	}
	msgResBytes, err := json.Marshal(msgRes)
	if err != nil {
		log.Printf("handleUserAuthenticate: Failed to marshal response into JSON: %v", err)
		closeConnection(c, websocket.CloseProtocolError, CloseReasonInternalServer)
		return
	}
	// Only broadcast message if successful, otherwise send only to the client
	if msgRes.Success == true {
		c.ws.boardHubs[params.BoardID].broadcast <- msgResBytes
	} else if msgRes.Success == false {
		c.send <- msgResBytes
	}
}

func handlePostCreate(c *Client, msgReq Request) {
	user := c.user
	if user.ID.String() == "" {
		closeConnection(c, websocket.ClosePolicyViolation, CloseReasonUnauthorized)
		return
	}
	var params ParamsPostCreate
	err := json.Unmarshal(msgReq.Params, &params)
	if err != nil {
		closeConnection(c, websocket.CloseInvalidFramePayloadData, CloseReasonBadParams)
		return
	}
	boardID := params.BoardID
	if !c.boards[boardID].canWrite {
		msgRes := buildErrorResponse(msgReq, ErrMsgUnauthorized)
		sendErrorMessage(c, msgRes)
		return
	}
	createPostInput := post.CreatePostInput{
		UserID:  user.ID.String(),
		BoardID: boardID,
		Content: params.Content,
		PosX:    params.PosX,
		PosY:    params.PosY,
		Color:   params.Color,
		Height:  params.Height,
		ZIndex:  params.ZIndex,
	}
	post, err := c.ws.postService.CreatePost(context.Background(), createPostInput)
	if err != nil {
		if errors.Is(err, validator.ValidationErrors{}) {
			validationErrMsg := v.GetValidationErrMsg(createPostInput, err)
			sendErrorMessage(c, buildErrorResponse(msgReq, validationErrMsg))
		} else {
			sendErrorMessage(c, buildErrorResponse(msgReq, ErrMsgInternalServer))
		}
		return
	}
	msgRes := ResponsePost{
		ResponseBase: ResponseBase{
			Event:   msgReq.Event,
			Success: true,
		},
		Result: post,
	}
	msgResBytes, err := json.Marshal(msgRes)

	if err != nil {
		log.Printf("handlePostCreate: Failed to marshal response into JSON: %v", err)
		closeConnection(c, websocket.CloseProtocolError, CloseReasonInternalServer)
		return
	}
	c.ws.boardHubs[boardID].broadcast <- msgResBytes
}

func handlePostFocus(c *Client, msgReq Request) {
	user := c.user
	if user.ID.String() == "" {
		closeConnection(c, websocket.ClosePolicyViolation, CloseReasonUnauthorized)
		return
	}
	var params ParamsPostFocus
	err := json.Unmarshal(msgReq.Params, &params)
	if err != nil {
		closeConnection(c, websocket.CloseInvalidFramePayloadData, CloseReasonBadParams)
		return
	}
	postID := params.ID
	boardID := params.BoardID
	if !c.boards[boardID].canWrite {
		msgRes := buildErrorResponse(msgReq, ErrMsgUnauthorized)
		sendErrorMessage(c, msgRes)
		return
	}
	msgRes := ResponsePostFocus{
		ResponseBase: ResponseBase{
			Event:   msgReq.Event,
			Success: true,
		},
		Result: ResultPostFocus{
			ID:      postID,
			BoardID: boardID,
			User:    *c.user,
		},
	}
	msgResBytes, err := json.Marshal(msgRes)

	if err != nil {
		log.Printf("handlePostFocus: Failed to marshal response into JSON: %v", err)
		closeConnection(c, websocket.CloseProtocolError, CloseReasonInternalServer)
		return
	}
	c.ws.boardHubs[boardID].broadcast <- msgResBytes
}

func handlePostUpdate(c *Client, msgReq Request) {
	user := c.user
	if user.ID.String() == "" {
		closeConnection(c, websocket.ClosePolicyViolation, CloseReasonUnauthorized)
		return
	}
	var params ParamsPostUpdate
	err := json.Unmarshal(msgReq.Params, &params)
	if err != nil {
		closeConnection(c, websocket.CloseInvalidFramePayloadData, CloseReasonBadParams)
		return
	}
	boardID := params.BoardID
	if !c.boards[boardID].canWrite {
		msgRes := buildErrorResponse(msgReq, ErrMsgUnauthorized)
		sendErrorMessage(c, msgRes)
		return
	}
	updatePostInput := post.UpdatePostInput{
		ID:      params.ID,
		Content: params.Content,
		PosX:    params.PosX,
		PosY:    params.PosY,
		Color:   params.Color,
		Height:  params.Height,
		ZIndex:  params.ZIndex,
	}
	post, err := c.ws.postService.UpdatePost(context.Background(), updatePostInput)
	if err != nil {
		if errors.Is(err, validator.ValidationErrors{}) {
			validationErrMsg := v.GetValidationErrMsg(updatePostInput, err)
			sendErrorMessage(c, buildErrorResponse(msgReq, validationErrMsg))
		} else {
			sendErrorMessage(c, buildErrorResponse(msgReq, ErrMsgInternalServer))
		}
		return
	}
	msgRes := ResponsePost{
		ResponseBase: ResponseBase{
			Event:   msgReq.Event,
			Success: true,
		},
		Result: post,
	}
	msgResBytes, err := json.Marshal(msgRes)

	if err != nil {
		log.Printf("handlePostUpdate: Failed to marshal response into JSON: %v", err)
		closeConnection(c, websocket.CloseProtocolError, CloseReasonInternalServer)
		return
	}
	c.ws.boardHubs[boardID].broadcast <- msgResBytes
}

func handlePostDelete(c *Client, msgReq Request) {
	user := c.user
	if user.ID.String() == "" {
		closeConnection(c, websocket.ClosePolicyViolation, CloseReasonUnauthorized)
		return
	}
	var params ParamsPostDelete
	err := json.Unmarshal(msgReq.Params, &params)
	if err != nil {
		closeConnection(c, websocket.CloseInvalidFramePayloadData, CloseReasonBadParams)
		return
	}
	postID := params.PostID
	boardID := params.BoardID
	if !c.boards[boardID].canWrite {
		msgRes := buildErrorResponse(msgReq, ErrMsgUnauthorized)
		sendErrorMessage(c, msgRes)
		return
	}
	err = c.ws.postService.DeletePost(context.Background(), postID)
	if err != nil {
		sendErrorMessage(c, buildErrorResponse(msgReq, ErrMsgInternalServer))
		return
	}
	msgRes := ResponsePostDelete{
		ResponseBase: ResponseBase{
			Event:   msgReq.Event,
			Success: true,
		},
		Result: ResultPostDelete{
			PostID:  postID,
			BoardID: boardID,
		},
	}
	msgResBytes, err := json.Marshal(msgRes)

	if err != nil {
		log.Printf("handlePostDelete: Failed to marshal response into JSON: %v", err)
		closeConnection(c, websocket.CloseProtocolError, CloseReasonInternalServer)
		return
	}
	c.ws.boardHubs[boardID].broadcast <- msgResBytes
}

func closeConnection(c *Client, statusCode int, text string) {
	c.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(statusCode, text), time.Now().Add(writeWait))
}
