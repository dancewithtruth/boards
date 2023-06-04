package ws2

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Wave-95/boards/server/internal/board"
	"github.com/Wave-95/boards/server/internal/jwt"
	"github.com/Wave-95/boards/server/internal/user"
	"github.com/Wave-95/boards/server/pkg/logger"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type WebSocket struct {
	userService  user.Service
	boardService board.Service
	jwtService   jwt.JWTService
	boardHubs    map[string]*Hub
	destroy      chan string
}

func NewWebSocket(userService user.Service, boardService board.Service, jwtService jwt.JWTService) *WebSocket {
	destroy := make(chan string)
	boardHubs := make(map[string]*Hub)
	go handleDestroy(destroy, boardHubs)
	return &WebSocket{
		userService:  userService,
		boardService: boardService,
		jwtService:   jwtService,
		boardHubs:    boardHubs,
		destroy:      destroy,
	}
}

func (ws *WebSocket) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
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

// handleConnectUser will authenticate the user and store a userId in the Client. It will also
// create a new board hub if it's the first connection to the provided board ID. It will register
// the client to the board and store the user's write permission for that board before returning a
// message response. The message response will contain the new user that was connected, along with the
// existing connected users
func handleConnectUser(c *Client, msgReq Request) {
	var params ParamsConnectUser
	err := json.Unmarshal(msgReq.Params, &params)
	if err != nil {
		//TODO: Close connection with close reason
		return
	}

	// Verify jwt token
	userId, err := c.ws.jwtService.VerifyToken(params.Jwt)
	if err != nil || userId == "" {
		//TODO: Close connection with close reason
		return
	}

	// Authorize user to board
	boardId := params.BoardId
	board, err := c.ws.boardService.GetBoardWithMembers(context.Background(), boardId)
	if err != nil {
		//TODO: Close connection with close reason
		return
	}
	if !hasAccess(board, userId) {
		//TODO: Close connection with close reason
	}

	// Create new board hub and run it if none exist
	if _, ok := c.ws.boardHubs[boardId]; !ok {
		c.ws.boardHubs[boardId] = newHub(boardId, c.ws.destroy)
		go c.ws.boardHubs[boardId].run()
	}
	hub := c.ws.boardHubs[boardId]
	// Store user ID, write permission, and hub location to client struct
	clientBoard := Board{
		hasWritePermission: true,
		hub:                hub,
	}
	c.userId = userId
	c.boards[boardId] = clientBoard
	// Register user to hub
	hub.register <- c

	// Broacast successful message response
	msgRes := ResponseConnectUser{
		Event:   EventConnectUser,
		Success: true,
		Result: ResultConnectUser{
			NewUser:       userId,
			ExistingUsers: hub.listConnectedUsers(),
			BoardId:       boardId,
		},
	}
	msgResBytes, err := json.Marshal(msgRes)
	if err != nil {
		//TODO: Handle error
	}
	hub.broadcast <- msgResBytes
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
