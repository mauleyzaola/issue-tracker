package tecsmtp

import (
	"fmt"
	"log"
)

type EmailConfig struct {
	Smtp    Smtp    `json:"smtp" bson:"smtp"`
	Mailgun Mailgun `json:"mailgun" bson:"mailgun"`
}

type Smtp struct {
	Username    string `json:"userName" bson:"userName"`
	Password    string `json:"password" bson:"password"`
	EmailServer string `json:"emailServer" bson:"emailServer"`
	Port        int    `json:"port" bson:"port"`
}

type Mailgun struct {
	Domain     string `json:"domain" bson:"domain"`
	PublicKey  string `json:"publicKey" bson:"publicKey"`
	PrivateKey string `json:"privateKey" bson:"privateKey"`
}

type EmailMessage struct {
	From        string   `json:"from"`
	To          string   `json:"to"`
	Tos         []string `json:"tos"`
	Subject     string   `json:"subject"`
	Body        string   `json:"body"`
	Attachments []string `json:"attachments"`
	BodyHtml    string   `json:"bodyHtml"`
}

func (t *EmailMessage) Validate() error {
	if len(t.From) == 0 {
		return fmt.Errorf("missing from parameter")
	} else if len(t.Subject) == 0 {
		return fmt.Errorf("missing subject")
	} else if len(t.To) == 0 && len(t.Tos) == 0 {
		return fmt.Errorf("missing to or tos parameter")
	}
	return nil
}

func (e *EmailConfig) IsSmtp() bool {
	return len(e.Mailgun.Domain) == 0 && len(e.Mailgun.PublicKey) == 0
}

func (e *EmailConfig) Validate() error {
	if e.IsSmtp() {
		if len(e.Smtp.EmailServer) == 0 {
			return fmt.Errorf("missing smtp server")
		} else if len(e.Smtp.Password) == 0 {
			return fmt.Errorf("missing smtp password")
		} else if len(e.Smtp.Username) == 0 {
			return fmt.Errorf("missing smtp user")
		} else if e.Smtp.Port <= 0 {
			return fmt.Errorf("missing port")
		}
	} else {
		if len(e.Mailgun.Domain) == 0 {
			return fmt.Errorf("missing domain")
		} else if len(e.Mailgun.PublicKey) == 0 {
			return fmt.Errorf("missing public domain")
		}
	}
	return nil
}

//Envia un mensaje por SMTP o Mailgun, segun este la configuracion
func SendMessage(config *EmailConfig, message *EmailMessage) (id string, response string, err error) {
	if config == nil {
		err = fmt.Errorf("missing config")
		return
	}

	err = config.Validate()
	if err != nil {
		return
	}

	err = message.Validate()
	if err != nil {
		return
	}

	if config.IsSmtp() {
		log.Println("Try to send SMTP email...")
		err = SendSmtp(config, message)
	} else {
		log.Println("Try to send Mailgun email...")
		response, id, err = SendMessageMailgun(config, message)
	}

	if err != nil {
		log.Println("Email error", err.Error())
	} else {
		log.Println("Email sent successfully")
	}
	return
}

//Valida que los caracteres en la direccion de correo sean ASCII
func ValidateEmailAddress(email string) bool {
	valid := len(email) != 0
	for i := range email {
		//z minuscula es el ultimo caracter valido
		if email[i] > 122 {
			valid = false
			break
		}
	}
	return valid
}
