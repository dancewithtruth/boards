package templates

import "os"

func BuildEmailVerification(to string, name string, code string) []byte {
	frontendURL := os.Getenv("FRONTEND_URL")
	link := frontendURL + "/verify-email?code=" + code
	msg := []byte("To: " + to + "\r\n" +
		"Subject: Boards: Verify your email address\r\n" +
		"\r\n" +
		"Hi " + name + ",\n\n" +
		"Please follow the link to verify your account: " + link + "\r\n")

	return msg
}

func BuildEmailInvite(to string, receiverName string, senderName string) []byte {
	frontendURL := os.Getenv("FRONTEND_URL")
	link := frontendURL + "/dashboard"
	msg := []byte("To: " + to + "\r\n" +
		"Subject: Boards: You've been invited to a board!\r\n" +
		"\r\n" +
		"Hi " + receiverName + ",\n\n" +
		"You've been invited to a board by " + senderName + ". You can accept or ignore the invitation from your dashboard: " + link + "\r\n")

	return msg
}
