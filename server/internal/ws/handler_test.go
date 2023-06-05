package ws

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wave-95/boards/server/internal/jwt"
	"github.com/Wave-95/boards/server/internal/models"
	"github.com/Wave-95/boards/server/internal/test"
	"github.com/Wave-95/boards/server/pkg/validator"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestHandleWebSocket(t *testing.T) {
	jwtService, server := setupServer(t)
	defer server.Close()

	user1, user1Conn := setupUserAndConnection(t, jwtService, server)
	defer user1Conn.Close()
	assertReceivedInitialMessage(t, user1Conn, user1, 0)

	user2, user2Conn := setupUserAndConnection(t, jwtService, server)
	defer user2Conn.Close()
	assertReceivedInitialMessage(t, user2Conn, user2, 1)
}

func setupServer(t *testing.T) (jwt.Service, *httptest.Server) {
	jwtService := jwt.New("jwt_secret", 1)
	validator := validator.New()
	wsAPI := NewAPI(jwtService, validator)
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", wsAPI.HandleWebSocket)
	server := httptest.NewServer(mux)
	return jwtService, server
}

func setupUserAndConnection(t *testing.T, jwtService jwt.Service, server *httptest.Server) (models.User, *websocket.Conn) {
	user := test.NewUser()
	token, err := jwtService.GenerateToken(user.Id.String())
	if err != nil {
		t.Fatalf("Failed to generate test jwt token: %v", err)
	}
	// Establish first WebSocket connection
	url := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws?jwt=" + token
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}
	return user, conn
}

func assertReceivedInitialMessage(t *testing.T, conn *websocket.Conn, user models.User, numExisting int) {
	// Assert first message received from WebSocket is a json with an array containing one user ID
	_, bytes, _ := conn.ReadMessage()
	var message ResponseUserConnect
	if err := json.Unmarshal(bytes, &message); err != nil {
		t.Fatalf("Failed to unmarshal first message: %v", err)
	}
	assert.Equal(t, EventBoardConnectUser, message.Event)
	assert.Equal(t, numExisting, len(message.Result.ExistingUsers))
	assert.Equal(t, user.Id.String(), message.Result.NewUser)
}
