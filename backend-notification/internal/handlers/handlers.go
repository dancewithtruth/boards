package handlers

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/Wave-95/boards/backend-notification/clients/boards"
	"github.com/Wave-95/boards/backend-notification/clients/email"
	"github.com/Wave-95/boards/backend-notification/constants/payloads"
	"github.com/Wave-95/boards/backend-notification/constants/queues"
	"github.com/Wave-95/boards/backend-notification/constants/tasks"
	"github.com/Wave-95/boards/backend-notification/internal/code"
	"github.com/Wave-95/boards/backend-notification/internal/templates"
	"github.com/Wave-95/boards/wrappers/amqp"
)

type TaskHandler struct {
	emailClient  *email.EmailClient
	boardsClient *boards.BoardsClient
	amqp         amqp.Amqp
}

func New(emailClient *email.EmailClient, boardsClient *boards.BoardsClient, amqp amqp.Amqp) TaskHandler {
	return TaskHandler{emailClient: emailClient, boardsClient: boardsClient, amqp: amqp}
}

func (th *TaskHandler) RegisterHandlers() {
	th.amqp.AddHandler(tasks.EmailInvite, th.emailInviteHandler)
	th.amqp.AddHandler(tasks.EmailVerification, th.emailVerificationHandler)
}

func (th *TaskHandler) Run() error {
	th.RegisterHandlers()
	return th.amqp.Consume(queues.Notification) // 10s
}

// emailVerificationHandler will create an email verification record and send the user
// a verification email containing the verification link
func (th *TaskHandler) emailVerificationHandler(payload []byte) error {
	var user payloads.User
	err := json.Unmarshal(payload, &user)
	if err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	verificationCode := code.Generate()
	reqBody := boards.CreateEmailVerificationInput{
		UserID: user.ID,
		Code:   verificationCode,
	}

	resp, err := th.boardsClient.CreateEmailVerification(reqBody)
	if err != nil {
		return fmt.Errorf("failed to prepare create email verification request: %w", err)
	}
	defer resp.Body.Close()

	// Send verification email
	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		emailTemplate := templates.BuildEmailVerification(user.Email, user.Name, verificationCode)
		err = th.emailClient.Send(user.Email, emailTemplate)
		if err != nil {
			return fmt.Errorf("failed to send verification email: %w", err)
		}
		return nil
	}

	return fmt.Errorf("failed to create email verification: %w", err)
}

func (th *TaskHandler) emailInviteHandler(payload []byte) error {
	var invite payloads.Invite
	err := json.Unmarshal(payload, &invite)
	if err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	resp, err := th.boardsClient.GetInvite(invite.ID)
	if err != nil {
		return fmt.Errorf("failed to get invite details: %w", err)
	}
	defer resp.Body.Close()
	inviteData, err := io.ReadAll(resp.Body)

	var inviteResponse InviteResponse
	err = json.Unmarshal(inviteData, &inviteResponse)
	if err != nil {
		return fmt.Errorf("failed to unmarshal invite details: %w", err)
	}

	emailBody := templates.BuildEmailInvite(inviteResponse.Receiver.Email, inviteResponse.Receiver.Name, inviteResponse.Sender.Name)
	return th.emailClient.Send(inviteResponse.Receiver.Email, emailBody)
}
