package ws

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Wave-95/boards/server/internal/endpoint"
	"github.com/Wave-95/boards/server/internal/post"
	"github.com/Wave-95/boards/server/pkg/logger"
	"github.com/Wave-95/boards/server/pkg/validator"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	ErrMsgInvalidToken = "Missing or invalid jwt token. Please ensure jwt is passed as query param to WebSocket endpoint"
)

// HandleWebSocket is the handler for the "/ws" endpoint and is responsible for authenticating the user
// and upgrading the http connection to a WebSocket. Once a connection is established, it registers the
// client to a hub, sends an initial message, and listens for incoming messages via client.ReadPump()
func (api *API) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logger.FromContext(ctx)
	// Parse jwt token from query param
	queryParams := r.URL.Query()
	token := queryParams.Get("jwt")

	// Verify jwt token
	userId, err := api.jwtService.VerifyToken(token)
	if err != nil || userId == "" {
		endpoint.WriteWithError(w, http.StatusUnauthorized, ErrMsgInvalidToken)
		return
	}

	// Upgrade connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Errorf("handler: failed to upgrade connection: %v", err)
		return
	}

	// Initialize and register new client
	client := &Client{
		api:    api,
		userId: userId,
		hub:    api.hub,
		conn:   conn,
		send:   make(chan []byte, 256),
	}
	existingUsers := api.hub.getUsers()
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()

	// Broadcast users connected
	if err := broadcastConnectUser(client, existingUsers); err != nil {
		logger.Errorf("handler: failed to broadcast user connected message: %v", err)
		client.hub.unregister <- client
	}

}

// HandleMessage handles every message that comes in. It's responsible for handling errors
// that occur when unmarshaling the []byte message into its relevant struct. HandleMessage
// will close the connection with a close reason for bad requests. If there's no formatting
// issue, the handler will either send a successful or unsuccessful message response.
func HandleMessage(c *Client, msg []byte) {
	// Identify message event
	var msgReq Request
	err := json.Unmarshal(msg, &msgReq)
	if err != nil {
		// TODO: Close connection with close reason
		// closeConnectionWithReason(statusCode int, errorMessage string)
		return
	}

	// Route message and handle accordingly
	switch msgReq.Event {
	case "":
		// TODO: Close connection with close reason
		return
	case EventPostCreate:
		handleCreatePost(c, msgReq)
	default:
		// TODO: Close connection with close reason
		return
	}
}

func handleCreatePost(c *Client, msgReq Request) {
	// Unmarshal input
	var input post.CreatePostInput
	err := json.Unmarshal(msgReq.Params, &input)
	if err != nil {
		// TODO: Close connnection with close reason
		return
	}
	// Validate create post payload
	err = c.api.validator.Struct(input)
	if err != nil {
		sendMsgValidationErr(c, msgReq.Event, input, err)
		return
	}

	// TODO: Create post
	// post, err := c.api.postService.CreatePost(payload)
	// if err != nil {
	// 	sendMessageInternalServerErr(c, ErrorMessageInternalCreatePost)
	// 	return
	// }
}

func sendMsgValidationErr(c *Client, event string, s interface{}, err error) {
	errMsg := "Invalid request"
	validationErrMsg := validator.GetValidationErrMsg(s, err)
	if validationErrMsg != "" {
		errMsg = validationErrMsg
	}
	c.send <- buildErrorResponse(event, errMsg)
}

func buildErrorResponse(event string, errMsg string) []byte {
	messageResponse := Response{
		Event:        event,
		ErrorMessage: errMsg,
	}
	messageResponseBytes, err := json.Marshal(messageResponse)
	if err != nil {
		log.Printf("Failed to marshal error response into json:%v", err)
	}
	return messageResponseBytes
}

func broadcastConnectUser(client *Client, existingUsers []string) error {
	messageUserConnected := ResponseUserConnect{
		Response: Response{
			Event:   EventBoardConnectUser,
			Success: true,
		},
		Result: ResultUserConnect{
			NewUser:       client.userId,
			ExistingUsers: existingUsers,
		},
	}
	messageUserConnectedBytes, err := json.Marshal(messageUserConnected)
	if err != nil {
		return err
	}
	client.hub.broadcast <- messageUserConnectedBytes
	return nil
}
