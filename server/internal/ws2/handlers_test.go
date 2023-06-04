package ws2

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wave-95/boards/server/internal/board"
	"github.com/Wave-95/boards/server/internal/jwt"
	"github.com/Wave-95/boards/server/internal/models"
	"github.com/Wave-95/boards/server/internal/test"
	"github.com/Wave-95/boards/server/internal/user"
	"github.com/Wave-95/boards/server/pkg/validator"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestHandleWebSocket(t *testing.T) {
	// Setup server
	// mux := http.NewServeMux()
	// mux.HandleFunc("/ws", HandleWebSocket)
	// server := httptest.NewServer(mux)
	// wsUrl := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	// t.Run("can establish a connection", func(t *testing.T) {
	// 	c, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	// 	assert.NoError(t, err)
	// 	assert.NotNil(t, c)
	// })

	t.Run("can send an initial board.connect_user request", func(t *testing.T) {
		// Create test user
		testUser := test.NewUser()
		mockUserRepo := user.NewMockRepository(make(map[uuid.UUID]models.User))
		err := mockUserRepo.CreateUser(context.Background(), testUser)
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		// Create test board
		mockBoardRepo := board.NewMockRepository(make(map[uuid.UUID]models.Board))
		testBoard := test.NewBoard(testUser.Id)
		err = mockBoardRepo.CreateBoard(context.Background(), testBoard)
		if err != nil {
			t.Fatalf("Failed to create test board: %v", err)
		}

		// Create test JWT token
		jwtService := jwt.New("jwt_secret", 1)
		token, err := jwtService.GenerateToken(testUser.Id.String())
		if err != nil {
			t.Fatalf("Failed to generate test JWT token: %v", err)
		}

		// Set up mock user and board service
		validator := validator.New()
		mockUserService := user.NewService(mockUserRepo, validator)
		mockBoardService := board.NewService(mockBoardRepo, validator)

		ws := NewWebSocket(mockUserService, mockBoardService, jwtService)
		mux := http.NewServeMux()
		mux.HandleFunc("/ws", ws.HandleWebSocket)
		server := httptest.NewServer(mux)
		wsUrl := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

		c, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
		if err != nil {
			t.Fatalf("Failed to establish connection: %v", err)
		}

		// Prepare initial message request
		msgReq := RequestConnectUser{
			Event: EventConnectUser,
			Params: ParamsConnectUser{
				Jwt:     token,
				BoardId: testBoard.Id.String(),
			},
		}

		// Write message to WebSocket
		if err := c.WriteJSON(msgReq); err != nil {
			t.Fatalf("Failed to write JSON for message request: %v", err)
		}

		// Read message from WebSocket and check fields
		_, msgRes, err := c.ReadMessage()
		var resConnectUser ResponseConnectUser
		json.Unmarshal(msgRes, &resConnectUser)
		assert.Equal(t, testUser.Id.String(), resConnectUser.Result.NewUser, "user ID from JWT does not match user ID returned in board.connect_user response")
	})
}
