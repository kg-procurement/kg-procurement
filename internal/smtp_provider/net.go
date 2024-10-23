//go:generate mockgen -typed -source=net.go -destination=net_mock.go -package=smtp_provider
package smtp_provider

import (
	"fmt"
	"kg/procurement/cmd/config"
	"log"
	gosmtp "net/smtp"
	"strings"
)

type SMTPClient interface {
	SendMail(addr string, a gosmtp.Auth, from string, to []string, msg []byte) error
}

type RealSMTPClient struct{}

func (c *RealSMTPClient) SendMail(addr string, a gosmtp.Auth, from string, to []string, msg []byte) error {
	return gosmtp.SendMail(addr, a, from, to, msg)
}

type nativeSMTP struct {
	cfg        config.SMTP
	auth       gosmtp.Auth
	smtpClient SMTPClient
}

func (n nativeSMTP) SendEmail(email Email) error {
	payloadByte := n.buildPayloadFromEmail(email)
	smtpAddr := fmt.Sprintf("%s:%s", n.cfg.Host, n.cfg.Port)

	err := n.smtpClient.SendMail(smtpAddr, n.auth, n.cfg.AuthEmail, append(email.To, email.CC...), payloadByte)
	if err != nil {
		log.Printf("smtp error sending email: %v", err)
		return err
	}

	return nil
}

func (n nativeSMTP) buildPayloadFromEmail(email Email) []byte {
	payload := "From: " + n.cfg.SenderName + "\n" +
		"To: " + strings.Join(email.To, ",") + "\n" +
		"Cc: " + strings.Join(email.CC, ",") + "\n" +
		"Subject: " + email.Subject + "\n\n" +
		email.Body

	return []byte(payload)
}

func NewNativeSMTP(cfg config.SMTP) *nativeSMTP {
	return &nativeSMTP{
		cfg:        cfg,
		auth:       gosmtp.PlainAuth("", cfg.AuthEmail, cfg.AuthPassword, cfg.Host),
		smtpClient: &RealSMTPClient{},
	}
}
