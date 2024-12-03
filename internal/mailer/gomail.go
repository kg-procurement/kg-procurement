//go:generate mockgen -typed -source=gomail.go -destination=gomail_mock.go -package=mailer
package mailer

import (
	"fmt"
	"io"
	"kg/procurement/cmd/config"
	"kg/procurement/cmd/utils"
	"strconv"

	"gopkg.in/gomail.v2"
)

// Dialer is an interface for sending emails.
type Dialer interface {
	DialAndSend(m ...*gomail.Message) error
}

type gomailSMTP struct {
	dialer Dialer
}

func (g gomailSMTP) SendEmail(email Email) error {
	message := g.buildEmailPayload(email)

	if err := g.dialer.DialAndSend(message); err != nil {
		utils.Logger.Errorf("Error sending email: %v", err)
		return err
	}

	return nil
}

func (gomailSMTP) buildEmailPayload(email Email) *gomail.Message {
	m := gomail.NewMessage()

	m.SetHeader("From", email.From)
	m.SetHeader("To", email.To...)
	m.SetHeader("Subject", email.Subject)
	m.SetHeader("Cc", email.CC...)
	m.SetBody("text/plain", email.Body)

	for _, attachment := range email.Attachments {
		m.Attach(attachment.Filename, gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(attachment.Data)
			return err
		}), gomail.SetHeader(map[string][]string{
			"Content-Type":        {attachment.MIMEType},
			"Content-Disposition": {fmt.Sprintf(`attachment; filename="%s"`, attachment.Filename)},
		}))
	}

	return m
}

func NewGomailSMTP(cfg config.SMTP) *gomailSMTP {
	portInt, err := strconv.Atoi(cfg.Port)
	if err != nil {
		panic(fmt.Errorf("invalid port number: %v", cfg.Port))
	}

	return &gomailSMTP{
		dialer: gomail.NewDialer(cfg.Host, portInt, cfg.AuthEmail, cfg.AuthPassword),
	}
}
