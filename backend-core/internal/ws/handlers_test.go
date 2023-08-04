package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wave-95/boards/backend-core/internal/board"
	"github.com/Wave-95/boards/backend-core/internal/config"
	"github.com/Wave-95/boards/backend-core/internal/jwt"
	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/Wave-95/boards/backend-core/internal/post"
	"github.com/Wave-95/boards/backend-core/internal/test"
	"github.com/Wave-95/boards/backend-core/internal/user"
	"github.com/Wave-95/boards/backend-core/pkg/validator"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestHandleWebSocket(t *testing.T) {
	// Set up mock repositories
	mockUserRepo := user.NewMockRepository()
	mockBoardRepo := board.NewMockRepository()
	mockPostRepo := post.NewMockRepository()

	// Set up mock services
	validator := validator.New()
	mockUserService := user.NewService(mockUserRepo, validator)
	mockBoardService := board.NewService(mockBoardRepo, validator)
	mockPostService := post.NewService(mockPostRepo)
	jwtService := jwt.New("jwt_secret", 1)

	// Set up server
	redisConfig := config.RedisConfig{Host: "redis-ws", Port: "6379"}
	ws := NewWebSocket(mockUserService, mockBoardService, mockPostService, jwtService, redisConfig)
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", ws.HandleConnection)
	server := httptest.NewServer(mux)
	t.Run("user.authenticate", func(t *testing.T) {
		t.Run("successfully authenticate user and return response", func(t *testing.T) {
			// Add test user to repo
			testUser := test.NewUser()
			err := mockUserRepo.CreateUser(context.Background(), testUser)
			if err != nil {
				t.Fatalf("Failed to create test user: %v", err)
			}
			c := setupConnection(t, server)
			// Generate test token
			token, err := jwtService.GenerateToken(testUser.ID.String())
			if err != nil {
				t.Fatalf("Failed to generate test JWT token: %v", err)
			}

			// Prepare initial message request
			msgReq := RequestUserAuthenticate{
				Event: EventUserAuthenticate,
				Params: ParamsUserAuthenticate{
					Jwt: token,
				},
			}

			// Write message to WebSocket
			if err := c.WriteJSON(msgReq); err != nil {
				t.Fatalf("Failed to write JSON for message request: %v", err)
			}

			// Read message from WebSocket and check fields
			_, msgRes, err := c.ReadMessage()
			if err != nil {
				assert.FailNow(t, "Failed to read message", err)
			}
			var resUserAuthenticate ResponseUserAuthenticate
			if err := json.Unmarshal(msgRes, &resUserAuthenticate); err != nil {
				assert.FailNow(t, "Failed to unmarshal JSON", err)
			}
			assert.Equal(t, true, resUserAuthenticate.Success, "expected user.authenticate response to be successful")
			assert.Equal(t, testUser.ID, resUserAuthenticate.Result.User.ID, "user ID from JWT does not match user ID returned in response")
		})

		t.Run("bad params result in connection close", func(t *testing.T) {
			c := setupConnection(t, server)
			requestWithBadParams := []byte(`{"event":"user.authenticate", "params": {"jwt": 1}}`)
			if err := c.WriteMessage(websocket.TextMessage, requestWithBadParams); err != nil {
				t.Fatalf("Failed to write JSON for message request: %v", err)
			}
			_, _, err := c.ReadMessage()
			closeErr := &websocket.CloseError{}
			assert.ErrorAs(t, err, &closeErr, "expected conn to be closed")
			assert.Equal(t, websocket.CloseInvalidFramePayloadData, closeErr.Code, "expected status code 1007")
			assert.Equal(t, CloseReasonBadParams, closeErr.Text)
		})

		t.Run("invalid JWT token", func(t *testing.T) {
			c := setupConnection(t, server)
			requestWithInvalidJwt := []byte(`{"event":"user.authenticate", "params": {"jwt": "invalidjwt"}}`)
			if err := c.WriteMessage(websocket.TextMessage, requestWithInvalidJwt); err != nil {
				t.Fatalf("Failed to write JSON for message request: %v", err)
			}
			_, msgRes, err := c.ReadMessage()
			var resUserAuthenticate ResponseUserAuthenticate
			if err := json.Unmarshal(msgRes, &resUserAuthenticate); err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}
			assert.NoError(t, err)
			assert.Equal(t, false, resUserAuthenticate.Success)
			assert.Equal(t, ErrMsgInvalidJwt, resUserAuthenticate.ErrorMessage)
		})
	})

	t.Run("board.connect", func(t *testing.T) {
		// Add test user to repo
		testUser := test.NewUser()
		err := mockUserRepo.CreateUser(context.Background(), testUser)
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		// // Add test board to repo
		testBoard := test.NewBoard(testUser.ID)
		err = mockBoardRepo.CreateBoard(context.Background(), testBoard)
		if err != nil {
			t.Fatalf("Failed to create test board: %v", err)
		}
		t.Run("successfully connect user to board and return response", func(t *testing.T) {
			c := setupConnection(t, server)
			authenticateUser(t, c, jwtService, testUser)

			// Prepare board connect request
			msgReq := RequestBoardConnect{
				Event: EventBoardConnect,
				Params: ParamsBoardConnect{
					BoardID: testBoard.ID.String(),
				},
			}
			// Send request
			if err := c.WriteJSON(msgReq); err != nil {
				t.Fatalf("Failed to write JSON for message request: %v", err)
			}
			_, msgRes, err := c.ReadMessage()
			if err != nil {
				t.Fatalf("Could not connect user to board: %v", err)
			}
			var resBoardConnect ResponseBoardConnect
			err = json.Unmarshal(msgRes, &resBoardConnect)
			if err != nil {
				t.Fatalf("Failed to unmarshal board connect response to Go struct: %v", err)
			}
			assert.Equal(t, testUser.ID, resBoardConnect.Result.NewUser.ID)
		})

		t.Run("user is not authenticated and cannot connect to board, close connection", func(t *testing.T) {
			c := setupConnection(t, server)
			msgReq := RequestBoardConnect{
				Event: EventBoardConnect,
				Params: ParamsBoardConnect{
					BoardID: testBoard.ID.String(),
				},
			}
			// Send request
			if err := c.WriteJSON(msgReq); err != nil {
				t.Fatalf("Failed to write JSON for message request: %v", err)
			}
			_, _, err := c.ReadMessage()
			closeErr := &websocket.CloseError{}
			assert.ErrorAs(t, err, &closeErr, "expected conn to be closed")
			assert.Equal(t, websocket.ClosePolicyViolation, closeErr.Code, "expected status code 1008")
			assert.Equal(t, CloseReasonUnauthorized, closeErr.Text)
		})

		t.Run("user is not authorized--should return an error response", func(t *testing.T) {
			c := setupConnection(t, server)
			authenticateUser(t, c, jwtService, testUser)
			// Prepare board connect request
			msgReq := RequestBoardConnect{
				Event: EventBoardConnect,
				Params: ParamsBoardConnect{
					BoardID: uuid.New().String(),
				},
			}
			// Send request
			if err := c.WriteJSON(msgReq); err != nil {
				t.Fatalf("Failed to write JSON for message request: %v", err)
			}
			_, msgRes, err := c.ReadMessage()
			if err != nil {
				t.Fatalf("Could not connect user to board: %v", err)
			}
			var resBoardConnect ResponseBoardConnect
			err = json.Unmarshal(msgRes, &resBoardConnect)
			if err != nil {
				t.Fatalf("Failed to unmarshal board connect response to Go struct: %v", err)
			}
			assert.Equal(t, false, resBoardConnect.Success)
			assert.Equal(t, ErrMsgBoardNotFound, resBoardConnect.ErrorMessage)
		})
	})
}

func setupConnection(t *testing.T, server *httptest.Server) *websocket.Conn {
	// Connect to WebSocket
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to establish connection: %v", err)
	}
	return c
}

func authenticateUser(t *testing.T, c *websocket.Conn, jwtService jwt.Service, testUser models.User) {
	// Generate test token
	token, err := jwtService.GenerateToken(testUser.ID.String())
	if err != nil {
		t.Fatalf("Failed to generate test JWT token: %v", err)
	}
	// Prepare authenticate request
	msgReq := RequestUserAuthenticate{
		Event: EventUserAuthenticate,
		Params: ParamsUserAuthenticate{
			Jwt: token,
		},
	}
	// Send authenticate request
	if err := c.WriteJSON(msgReq); err != nil {
		t.Fatalf("Failed to write JSON for message request: %v", err)
	}
	_, _, err = c.ReadMessage()
	if err != nil {
		t.Fatalf("Could not authenticate user: %v", err)
	}
}
