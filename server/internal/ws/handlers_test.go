package ws

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
	"github.com/Wave-95/boards/server/internal/post"
	"github.com/Wave-95/boards/server/internal/test"
	"github.com/Wave-95/boards/server/internal/user"
	"github.com/Wave-95/boards/server/pkg/validator"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestHandleWebSocket(t *testing.T) {
	// Set up mock repositories
	mockUserRepo := user.NewMockRepository(make(map[uuid.UUID]models.User))
	mockBoardRepo := board.NewMockRepository(make(map[uuid.UUID]models.Board))
	mockPostRepo := post.NewMockRepository()

	// Set up mock services
	validator := validator.New()
	mockUserService := user.NewService(mockUserRepo, validator)
	mockBoardService := board.NewService(mockBoardRepo, validator)
	mockPostService := post.NewService(mockPostRepo)
	jwtService := jwt.New("jwt_secret", 1)

	// Set up server
	ws := NewWebSocket(mockUserService, mockBoardService, mockPostService, jwtService)
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
			token, err := jwtService.GenerateToken(testUser.Id.String())
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
			var resUserAuthenticate ResponseUserAuthenticate
			json.Unmarshal(msgRes, &resUserAuthenticate)
			assert.Equal(t, true, resUserAuthenticate.Success, "expected user.authenticate response to be successful")
			assert.Equal(t, testUser.Id, resUserAuthenticate.Result.User.Id, "user ID from JWT does not match user ID returned in response")
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
			json.Unmarshal(msgRes, &resUserAuthenticate)
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
		testBoard := test.NewBoard(testUser.Id)
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
					BoardId: testBoard.Id.String(),
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
			assert.Equal(t, testUser.Id.String(), resBoardConnect.Result.UserId)
		})

		t.Run("user is not authenticated and cannot connect to board, close connection", func(t *testing.T) {
			c := setupConnection(t, server)
			msgReq := RequestBoardConnect{
				Event: EventBoardConnect,
				Params: ParamsBoardConnect{
					BoardId: testBoard.Id.String(),
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
			randTestUser := test.NewUser()
			c := setupConnection(t, server)
			authenticateUser(t, c, jwtService, randTestUser)

			// Prepare board connect request
			msgReq := RequestBoardConnect{
				Event: EventBoardConnect,
				Params: ParamsBoardConnect{
					BoardId: testBoard.Id.String(),
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

func setupServer(t *testing.T, testUser models.User, testBoard models.Board, jwtService jwt.Service) *httptest.Server {
	// Set up mock user repo
	mockUserRepo := user.NewMockRepository(make(map[uuid.UUID]models.User))
	err := mockUserRepo.CreateUser(context.Background(), testUser)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	// Set up mock board repo
	mockBoardRepo := board.NewMockRepository(make(map[uuid.UUID]models.Board))
	err = mockBoardRepo.CreateBoard(context.Background(), testBoard)
	if err != nil {
		t.Fatalf("Failed to create test board: %v", err)
	}
	// Set up mock board repo
	mockPostRepo := post.NewMockRepository()

	// Set up mock user and board service
	validator := validator.New()
	mockUserService := user.NewService(mockUserRepo, validator)
	mockBoardService := board.NewService(mockBoardRepo, validator)
	mockPostService := post.NewService(mockPostRepo)

	// Set up server
	ws := NewWebSocket(mockUserService, mockBoardService, mockPostService, jwtService)
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", ws.HandleConnection)
	return httptest.NewServer(mux)
}

func setupConnection(t *testing.T, server *httptest.Server) *websocket.Conn {
	// Connect to WebSocket
	wsUrl := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	c, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		t.Fatalf("Failed to establish connection: %v", err)
	}
	return c
}

func authenticateUser(t *testing.T, c *websocket.Conn, jwtService jwt.Service, testUser models.User) {
	// Generate test token
	token, err := jwtService.GenerateToken(testUser.Id.String())
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
