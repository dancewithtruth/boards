package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"

	"github.com/Wave-95/boards/backend-notification/constants/payloads"
	"github.com/Wave-95/boards/backend-notification/constants/queues"
	"github.com/Wave-95/boards/backend-notification/constants/tasks"
	"github.com/Wave-95/boards/backend-notification/internal/code"
	"github.com/Wave-95/boards/backend-notification/internal/email"
	"github.com/Wave-95/boards/wrappers/amqp"
)

type TaskHandler struct {
	emailClient *email.EmailClient
	amqp        amqp.Amqp
}

func New(emailClient *email.EmailClient, amqp amqp.Amqp) TaskHandler {
	return TaskHandler{emailClient: emailClient, amqp: amqp}
}

func (th *TaskHandler) RegisterHandlers() {
	th.amqp.AddHandler(tasks.EmailInvite, th.emailInviteHandler)
	th.amqp.AddHandler(tasks.EmailVerification, th.emailVerificationHandler)
}

func (th *TaskHandler) Run() error {
	th.RegisterHandlers()
	return th.amqp.Consume(queues.Notification)
}

// TODO: Use env vars for URLs

// emailVerificationHandler will create an email verification record and send the user
// a verification email containing the verification link
func (th *TaskHandler) emailVerificationHandler(payload []byte) error {
	var user payloads.User
	err := json.Unmarshal(payload, &user)
	if err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	verificationCode := code.Generate()
	reqBody := CreateEmailVerificationInput{
		UserID: user.ID,
		Code:   verificationCode,
	}
	data, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body for creating email verification: %w", err)
	}

	// Create a new HTTP request with the POST method and data
	req, err := http.NewRequest("POST", "http://backend-core:8080/email-verifications", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to prepare create email verification request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Create a new HTTP client and send the request
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed create email verification API request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		//Send email
		msg := email.BuildEmailVerificationTemplate(user.Email)
		return th.emailClient.Send(user.Email, msg)
	}

	return fmt.Errorf("failed to create email verification: %w", err)
}

func (th *TaskHandler) emailInviteHandler(payload []byte) error {
	var invite payloads.Invite
	err := json.Unmarshal(payload, &invite)
	if err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	httpClient := http.DefaultClient
	resp, err := httpClient.Get("http://backend-core:8080/invites/" + invite.ID)
	if err != nil {
		// handle error
		fmt.Println("Encountered error")
	}
	defer resp.Body.Close()

	toList := []string{"wuvictor95@gmail.com"}

	host := "smtp.gmail.com"

	// Its the default port of smtp server
	port := "587"

	// This is the message to send in the mail
	msg := "Hello geeks!!!"

	// We can't send strings directly in mail,
	// strings need to be converted into slice bytes
	emailBody := []byte(msg)

	// Set the auth details
	auth := smtp.PlainAuth("", "useboards@gmail.com", "wdlrmviaiwalxkkq", "smtp.gmail.com")

	err = smtp.SendMail(host+":"+port, auth, "useboards@gmail.com", toList, emailBody)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully sent mail to all user in toList")

	return nil
}
