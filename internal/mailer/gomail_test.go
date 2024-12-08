package mailer

import (
	"bytes"
	"encoding/base64"
	"errors"
	"kg/procurement/cmd/config"
	"testing"

	"github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func Test_NewGomailSMTP(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		_ = NewGomailSMTP(config.SMTP{Port: "587"})
	})

	t.Run("panic invalid port number", func(t *testing.T) {
		g := gomega.NewWithT(t)
		cfg := config.SMTP{
			Host:         "smtp.example.com",
			Port:         "oh no wat is this port",
			AuthEmail:    "sender@example.com",
			AuthPassword: "password",
			SenderName:   "Test Sender",
		}

		g.Expect(func() {
			NewGomailSMTP(cfg)
		}).To(gomega.Panic())
	})
}

func Test_buildEmailPayload(t *testing.T) {
	g := gomega.NewWithT(t)

	email := Email{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		CC:      []string{"cc@example.com"},
		Subject: "Test Subject",
		Body:    "This is the email body.",
		Attachments: []Attachment{
			{
				Filename: "test.txt",
				Data:     []byte("This is test attachment data."),
				MIMEType: "text/plain",
			},
		},
	}

	gSMTP := gomailSMTP{}
	res := gSMTP.buildEmailPayload(email)

	g.Expect(res.GetHeader("From")[0]).To(gomega.Equal(email.From))
	g.Expect(res.GetHeader("To")[0]).To(gomega.Equal(email.To[0]))
	g.Expect(res.GetHeader("Subject")[0]).To(gomega.Equal(email.Subject))
	g.Expect(res.GetHeader("Cc")[0]).To(gomega.Equal(email.CC[0]))

	var buf bytes.Buffer
	_, err := res.WriteTo(&buf)
	g.Expect(err).To(gomega.BeNil())

	mimeContent := buf.String()
	g.Expect(mimeContent).To(gomega.ContainSubstring("Content-Disposition: attachment; filename=\"test.txt\""))

	expectedAttachmentData := base64.StdEncoding.EncodeToString([]byte("This is test attachment data."))
	g.Expect(mimeContent).To(gomega.ContainSubstring(expectedAttachmentData))
}

func Test_GomailSendEmail(t *testing.T) {
	var (
		email = Email{
			To:      []string{},
			CC:      []string{"cc@example.com"},
			Subject: "help",
			Body:    "ikimashoo",
		}
	)

	t.Run("success", func(t *testing.T) {
		g := gomega.NewWithT(t)
		ctrl := gomock.NewController(t)
		mockDialer := NewMockDialer(ctrl)

		gm := &gomailSMTP{
			dialer: mockDialer,
		}

		mockDialer.EXPECT().DialAndSend(gomock.Any()).
			Times(1).
			Return(nil)

		err := gm.SendEmail(email)
		g.Expect(err).To(gomega.BeNil())
	})

	t.Run("success", func(t *testing.T) {
		g := gomega.NewWithT(t)
		ctrl := gomock.NewController(t)
		mockDialer := NewMockDialer(ctrl)

		gm := &gomailSMTP{
			dialer: mockDialer,
		}

		mockDialer.EXPECT().DialAndSend(gomock.Any()).
			Times(1).
			Return(errors.New("oh no"))

		err := gm.SendEmail(email)
		g.Expect(err).ToNot(gomega.BeNil())
	})
}
