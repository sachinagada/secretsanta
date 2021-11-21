package email

import (
	"os"
	"testing"
)

func TestSendMail(t *testing.T) {
	username := os.Getenv("SECRET_SANTA_USERNAME")
	password := os.Getenv("SECRET_SANTA_PASSWORD")

	// don't run this in github where it doesn't have the username and password
	// and also don't want to recieve an email everytime the tests are run
	if username == "" || password == "" {
		return
	}

	c := Config{
		Username: username,
		Password: password,
		SMTPHost: "smtp.gmail.com",
		SMPTPort: "587",
		Subject:  "Test Secret Santa",
	}

	m, mErr := NewMail(&c)
	if mErr != nil {
		t.Fatalf("unexpected error initializing mail: %s", mErr)
	}

	s := Santa{
		Name:      "TestName",
		Addr:      username,
		Recipient: "Test Recipient",
	}

	sendErr := m.SendEmail([]Santa{s})
	if sendErr != nil {
		t.Fatalf("unexpected error sending mail: %s", sendErr)
	}
}
