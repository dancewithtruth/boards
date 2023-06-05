package ws2

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Wave-95/boards/server/internal/board"
	"github.com/Wave-95/boards/server/pkg/logger"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
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
		boardId := <-destroy
		delete(boardHubs, boardId)
	}
}

// handleUserAuthenticate will authenticate the user and store a userId in the Client.
func handleUserAuthenticate(c *Client, msgReq Request) {
	// Unmarshal params
	var params ParamsUserAuthenticate
	err := json.Unmarshal(msgReq.Params, &params)
	if err != nil {
		closeConnection(c, websocket.CloseInvalidFramePayloadData, CloseReasonBadParams)
		return
	}
	// Verify jwt token
	userId, err := c.ws.jwtService.VerifyToken(params.Jwt)
	var msgRes ResponseUserAuthenticate
	if err != nil || userId == "" {
		// Prepare error response
		msgRes = ResponseUserAuthenticate{
			Event:        EventBoardConnect,
			Success:      false,
			ErrorMessage: ErrMsgInvalidJwt,
		}
	} else {
		// Prepare success response
		msgRes = ResponseUserAuthenticate{
			Event:   EventBoardConnect,
			Success: true,
			Result: ResultUserAuthenticate{
				UserId: userId,
			},
		}
		// Assign userId to client struct
		c.userId = userId
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
	userId := c.userId
	if userId == "" {
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
	boardId := params.BoardId
	board, err := c.ws.boardService.GetBoardWithMembers(context.Background(), boardId)
	var msgRes ResponseBoardConnect
	// If no access, return error response
	if err != nil || !hasAccess(board, userId) {
		msgRes = ResponseBoardConnect{
			Event:        EventBoardConnect,
			Success:      false,
			ErrorMessage: ErrMsgBoardNotFound,
		}
	} else {
		// If board does not exist as a hub, create one and run it
		if _, ok := c.ws.boardHubs[boardId]; !ok {
			c.ws.boardHubs[boardId] = newHub(boardId, c.ws.destroy)
			go c.ws.boardHubs[boardId].run()
		}
		// Store write permission on the client's boards map
		boardHub := c.ws.boardHubs[boardId]
		c.boards[boardId] = Board{
			canWrite: true,
		}
		existingUsers := boardHub.listConnectedUsers()
		boardHub.register <- c

		// Broacast successful message response
		msgRes = ResponseBoardConnect{
			Event:   EventBoardConnect,
			Success: true,
			Result: ResultBoardConnect{
				BoardId:       boardId,
				UserId:        userId,
				ExistingUsers: existingUsers,
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
		c.ws.boardHubs[params.BoardId].broadcast <- msgResBytes
	} else if msgRes.Success == false {
		c.send <- msgResBytes
	}
}

func hasAccess(board board.BoardWithMembersDTO, userId string) bool {
	if board.UserId.String() == userId {
		return true
	}

	for _, member := range board.Members {
		if member.Id.String() == userId {
			return true
		}
	}
	return false
}

func closeConnection(c *Client, statusCode int, text string) {
	c.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(statusCode, text), time.Now().Add(writeWait))
}
