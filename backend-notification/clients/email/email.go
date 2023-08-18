package email

import (
	"net/smtp"
)

type EmailClient struct {
	from string
	auth smtp.Auth
	addr string
}

func NewClient(from, password, host, port string) *EmailClient {
	// Load env vars
	auth := smtp.PlainAuth("", from, password, host)
	addr := host + ":" + port
	return &EmailClient{from: from, auth: auth, addr: addr}
}

func (c *EmailClient) Send(to string, body []byte) error {
	toList := []string{to}
	return smtp.SendMail(c.addr, c.auth, c.from, toList, body)
}
