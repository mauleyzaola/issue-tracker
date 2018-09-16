package tecsmtp

import (
	"testing"
)

func TestEmailValidation(t *testing.T) {
	t.Log("Given an email configuration validate it has the required values")
	config := &EmailConfig{}
	config.Smtp.EmailServer = "smtp.google.com"
	err := config.Validate()
	if err == nil {
		t.Error("not equal")
		return
	}

	config.Smtp.Username = "username"
	err = config.Validate()
	if err == nil {
		t.Error("not equal")
		return
	}

	t.Log("Validate SMTP properties")
	config.Smtp.Password = "password"
	config.Smtp.Port = 25
	err = config.Validate()
	if !config.IsSmtp() {
		t.Error("not equal")
		return
	}

	t.Log("Validate Mailgun properties")
	config.Mailgun.Domain = "smtp.google.com"
	err = config.Validate()
	if err == nil {
		t.Error("not equal")
		return
	}

	config.Mailgun.PublicKey = "key"
	err = config.Validate()
	if err != nil {
		return
	}
	if config.IsSmtp() {
		t.Error("not equal")
		return
	}
}

func TestMessageValidation(t *testing.T) {
	t.Log("Given a new message, validate its properties are working as expected")

	msg := &EmailMessage{}
	err := msg.Validate()
	if err == nil {
		t.Error("not equal")
		return
	}
	t.Log(err)

	msg.From = "from@address.com"
	err = msg.Validate()
	if err == nil {
		t.Error("not equal")
		return
	}
	t.Log(err)

	msg.Subject = "message body"
	err = msg.Validate()
	if err == nil {
		t.Error("not equal")
		return
	}
	t.Log(err)

	msg.To = "recipient@gmail.com"
	err = msg.Validate()
	if err != nil {
		t.Error(err)
		return
	}
}
