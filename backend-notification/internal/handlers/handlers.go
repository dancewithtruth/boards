package handlers

import (
	"github.com/Wave-95/boards/backend-notification/constants/tasks"
	"github.com/Wave-95/boards/wrappers/amqp"
)

func Register(amqp amqp.Amqp) {
	amqp.AddHandler(tasks.EmailInvite, emailInviteHandler)
}

func emailInviteHandler(payload interface{}) error {
	//Implementation to fetch email info, build email template, and send email
	return nil
}
