package tecsmtp

import (
	"encoding/base64"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"
)

func encodeRFC2047(value string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{Name: value}
	return strings.Trim(addr.String(), " <>")
}

func SendSmtp(config *EmailConfig, message *EmailMessage) error {
	err := message.Validate()
	if err != nil {
		return err
	}

	err = config.Validate()
	if err != nil {
		return err
	}

	if !config.IsSmtp() {
		return fmt.Errorf("only smtp configuration is supported for smtp messages")
	}

	if len(message.Body) == 0 {
		return fmt.Errorf("html is not supported in smtp messages")
	}

	auth := smtp.PlainAuth(
		"",
		config.Smtp.Username,
		config.Smtp.Password,
		config.Smtp.EmailServer,
	)

	if len(message.To) != 0 {
		message.Tos = append(message.Tos, message.To)
	}

	if len(message.Tos) > 1 {
		return fmt.Errorf("only one recipient is supported in smtp messages")
	}

	from := mail.Address{Address: message.From}
	to := mail.Address{Address: message.Tos[0]}
	title := message.Subject
	body := message.Body

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = encodeRFC2047(title)
	header["MIME-Version"] = "1.0"

	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	msg := ""
	for k, v := range header {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	msg += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	return smtp.SendMail(
		fmt.Sprintf("%s:%v", config.Smtp.EmailServer, config.Smtp.Port),
		auth,
		from.Address,
		[]string{to.Address},
		[]byte(msg),
	)
}
