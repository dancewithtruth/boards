package ws

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wave-95/boards/server/internal/jwt"
	"github.com/Wave-95/boards/server/pkg/validator"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestHandleWebSocket(t *testing.T) {
	// Set up server
	jwtService := jwt.New("jwt_secret", 1)
	validator := validator.New()
	wsAPI := NewAPI(jwtService, validator)
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", wsAPI.HandleWebSocket)
	server := httptest.NewServer(mux)

	// Generate JWT tokens
	user1Id := uuid.New().String()
	user2Id := uuid.New().String()
	user1Token, err := jwtService.GenerateToken(user1Id)
	user2Token, err := jwtService.GenerateToken(user2Id)
	if err != nil {
		t.Fatalf("Failed to generate test jwt token: %v", err)
	}

	// Establish first WebSocket connection
	user1Url := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws?jwt=" + user1Token
	user1Conn, _, err := websocket.DefaultDialer.Dial(user1Url, nil)
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}
	defer user1Conn.Close()

	// Assert first message received from WebSocket is a json with an array containing one user ID
	_, connectUser1Bytes, _ := user1Conn.ReadMessage()
	var connectUser1Message MessageResponseConnectUser
	json.Unmarshal(connectUser1Bytes, &connectUser1Message)
	if len(connectUser1Message.Data.ExistingUsers) == 0 {
		t.Fatal("Failed to return user ID in first message")
	}
	assert.Equal(t, TypeConnectUser, connectUser1Message.Type)
	assert.Equal(t, 0, len(connectUser1Message.Data.ExistingUsers))
	assert.Equal(t, user1Id, connectUser1Message.Data.NewUser)

	// Establish a second WebSocket connection
	user2Url := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws?jwt=" + user2Token
	user2Conn, _, err := websocket.DefaultDialer.Dial(user2Url, nil)
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}
	defer user2Conn.Close()

	// Assert first message received from WebSocket is a json with an array containing two user IDs
	_, connectUser2Bytes, _ := user2Conn.ReadMessage()
	var connectUser2Message MessageResponseConnectUser
	json.Unmarshal(connectUser2Bytes, &connectUser2Message)
	if len(connectUser2Message.Data.ExistingUsers) == 0 {
		t.Fatal("Failed to return user ID in first message for client 2")
	}

	assert.Equal(t, 1, len(connectUser2Message.Data.ExistingUsers), "Expected array of length 1 for users connected")
	assert.Equal(t, user2Id, connectUser2Message.Data.NewUser)
}
