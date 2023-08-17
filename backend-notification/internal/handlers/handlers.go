package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/smtp"

	"github.com/Wave-95/boards/backend-notification/constants/payloads"
	"github.com/Wave-95/boards/backend-notification/constants/tasks"
	"github.com/Wave-95/boards/wrappers/amqp"
)

func Register(amqp amqp.Amqp) {
	amqp.AddHandler(tasks.EmailInvite, emailInviteHandler)
}

func emailInviteHandler(payload []byte) error {
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
	body, err := io.ReadAll(resp.Body)
	//Implementation to fetch email info, build email template, and send email
	fmt.Println(string(body))

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
