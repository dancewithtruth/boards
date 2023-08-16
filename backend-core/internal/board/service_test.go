package board

import (
	"context"
	"testing"

	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/Wave-95/boards/backend-core/internal/test"
	"github.com/Wave-95/boards/backend-core/pkg/validator"
	"github.com/Wave-95/boards/wrappers/amqp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	validator := validator.New()
	testUser := test.NewUser()
	mockBoardRepo := NewMockRepository()
	mockBoardRepo.AddUser(testUser)
	mockAmqp := amqp.NewMock()
	boardService := NewService(mockBoardRepo, mockAmqp, validator)
	assert.NotNil(t, boardService)
	t.Run("Create board", func(t *testing.T) {
		t.Run("without name or description", func(t *testing.T) {
			input := CreateBoardInput{
				UserID: testUser.ID.String(),
			}
			board, err := boardService.CreateBoard(context.Background(), input)
			assert.NoError(t, err)
			assert.Equal(t, "Board #1", *board.Name)
			assert.Equal(t, defaultBoardDescription, *board.Description)
		})

		t.Run("with name or description", func(t *testing.T) {
			customBoardName := "Custom Board Name"
			customBoardDescription := "Custom board description"
			input := CreateBoardInput{
				UserID:      testUser.ID.String(),
				Name:        &customBoardName,
				Description: &customBoardDescription,
			}
			board, err := boardService.CreateBoard(context.Background(), input)
			assert.NoError(t, err)
			assert.Equal(t, customBoardName, *board.Name)
			assert.Equal(t, customBoardDescription, *board.Description)
		})
	})

	t.Run("Get board", func(t *testing.T) {
		board, ok := getFirstBoard(mockBoardRepo.boards)
		if !ok {
			assert.FailNow(t, "expected a board to exist but got none")
		}
		board, err := boardService.GetBoard(context.Background(), board.ID.String())
		assert.NoError(t, err)
		assert.NotNil(t, board)
	})

	t.Run("List owned boards", func(t *testing.T) {
		boards, err := boardService.ListOwnedBoards(context.Background(), testUser.ID.String())
		assert.NoError(t, err)
		assert.Greater(t, len(boards), 0)
	})

	t.Run("Create board invites", func(t *testing.T) {
		// Setup test receivers
		receiver1 := test.NewUser()
		receiver2 := test.NewUser()

		// Setup test board
		createBoardInput := CreateBoardInput{
			UserID: testUser.ID.String(),
		}
		board, err := boardService.CreateBoard(context.Background(), createBoardInput)
		if err != nil {
			assert.FailNow(t, "Failed to create test board")
		}

		// Prepare board invites payload
		createBoardInvitesInput := CreateInvitesInput{
			BoardID:  board.ID.String(),
			SenderID: testUser.ID.String(),
			Invites: []struct {
				ReceiverID string `json:"receiver_id"`
			}{{ReceiverID: receiver1.ID.String()}, {ReceiverID: receiver2.ID.String()}},
		}
		invites, err := boardService.CreateInvites(context.Background(), createBoardInvitesInput)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(invites), "Expected an invites slice of length 2 to be returned, got ", len(invites))
	})
}

func getFirstBoard(m map[uuid.UUID]models.Board) (models.Board, bool) {
	for _, board := range m {
		return board, true
	}
	return models.Board{}, false
}
