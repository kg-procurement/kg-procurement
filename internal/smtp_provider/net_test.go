package smtp_provider

import (
	"errors"
	"fmt"
	"kg/procurement/cmd/config"
	"net/smtp"
	"strings"
	"testing"

	"github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func Test_NewNativeSMTP(t *testing.T) {
	_ = NewNativeSMTP(config.SMTP{})
}

func Test_buildPayloadFromEmail(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	nSMTP := NewNativeSMTP(config.SMTP{SenderName: "dlwlrma"})

	email := Email{
		To:      []string{},
		CC:      []string{"cc@example.com"},
		Subject: "help",
		Body:    "ikimashoo",
	}

	payload := "From: " + "dlwlrma" + "\n" +
		"To: " + strings.Join(email.To, ",") + "\n" +
		"Cc: " + strings.Join(email.CC, ",") + "\n" +
		"Subject: " + email.Subject + "\n\n" +
		email.Body

	payloadByte := nSMTP.buildPayloadFromEmail(email)

	g.Expect(string(payloadByte)).To(gomega.Equal(payload))
}

func Test_SendEmail(t *testing.T) {
	var (
		email = Email{
			To:      []string{},
			CC:      []string{"cc@example.com"},
			Subject: "help",
			Body:    "ikimashoo",
		}
		cfg = config.SMTP{
			Host:         "smtp.example.com",
			Port:         "587",
			AuthEmail:    "sender@example.com",
			AuthPassword: "password",
			SenderName:   "Test Sender",
		}
		smtpAddr = fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	)

	t.Run("success", func(t *testing.T) {
		g := gomega.NewGomegaWithT(t)
		ctrl := gomock.NewController(t)
		mockSMTPClient := NewMockSMTPClient(ctrl)

		n := &nativeSMTP{
			cfg:        cfg,
			auth:       smtp.PlainAuth("", cfg.AuthEmail, cfg.AuthPassword, cfg.Host),
			smtpClient: mockSMTPClient,
		}

		receivers := append(email.To, email.CC...)

		mockSMTPClient.EXPECT().
			SendMail(smtpAddr, n.auth, cfg.AuthEmail, receivers, n.buildPayloadFromEmail(email)).
			Return(nil)

		err := n.SendEmail(email)
		g.Expect(err).To(gomega.BeNil())
	})

	t.Run("error", func(t *testing.T) {
		g := gomega.NewGomegaWithT(t)
		ctrl := gomock.NewController(t)
		mockSMTPClient := NewMockSMTPClient(ctrl)

		n := &nativeSMTP{
			cfg:        cfg,
			auth:       smtp.PlainAuth("", cfg.AuthEmail, cfg.AuthPassword, cfg.Host),
			smtpClient: mockSMTPClient,
		}

		receivers := append(email.To, email.CC...)

		mockSMTPClient.EXPECT().
			SendMail(smtpAddr, n.auth, cfg.AuthEmail, receivers, n.buildPayloadFromEmail(email)).
			Return(errors.New("error"))

		err := n.SendEmail(email)
		g.Expect(err).ToNot(gomega.BeNil())
	})
}
