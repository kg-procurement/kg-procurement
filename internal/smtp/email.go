//go:generate mockgen -typed -source=email.go -destination=email_mock.go -package=smtp
package smtp

type Email struct {
	From    string
	To      []string
	CC      []string
	Subject string
	Body    string
}

type EmailSender interface {
	SendEmail(email Email) error
}
