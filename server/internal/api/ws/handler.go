package ws

import (
	"encoding/json"
	"net/http"

	"github.com/Wave-95/boards/server/internal/endpoint"
	"github.com/Wave-95/boards/server/pkg/logger"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	ErrMsgInvalidToken = "Missing or invalid jwt token. Please ensure jwt is passed as query param to WebSocket endpoint"
)

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

func broadcastConnectUser(client *Client, existingUsers []string) error {
	messageUserConnected := MessageResponseConnectUser{
		MessageResponse: MessageResponse{
			Type:   TypeConnectUser,
			Sender: client.userId,
		},
		Data: DataConnectUser{
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
