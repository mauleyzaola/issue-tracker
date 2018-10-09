package tecsmtp

import (
	"github.com/mailgun/mailgun-go"
)

func SendMessageMailgun(config *EmailConfig, message *EmailMessage) (id string, response string, err error) {
	gun := mailgun.NewMailgun(config.Mailgun.Domain, config.Mailgun.PublicKey)
	if len(message.BodyHtml) != 0 {
		message.Body = message.BodyHtml
	}

	var m *mailgun.Message
	if len(message.Tos) != 0 {
		m = mailgun.NewMessage(message.From, message.Subject, message.Body, message.Tos...)
	} else {
		m = mailgun.NewMessage(message.From, message.Subject, message.Body, message.To)
	}

	for _, i := range message.Attachments {
		m.AddAttachment(i)
	}

	if len(message.BodyHtml) != 0 {
		m.SetHtml(message.BodyHtml)
	}

	id, response, err = gun.Send(m)
	return
}
