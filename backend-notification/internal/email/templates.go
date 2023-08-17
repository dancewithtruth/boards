package email

func BuildEmailVerificationTemplate(to string) []byte {
	link := "http://useboards.com/verify-email"
	msg := []byte("To: " + to + "\r\n" +
		"Subject: Boards: Verify your email address\r\n" +
		"\r\n" +
		"Click the link to verify your email:" + link + "\r\n")

	return msg
}
