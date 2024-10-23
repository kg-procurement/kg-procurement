//go:generate mockgen -typed -source=email.go -destination=email_mock.go -package=smtp_provider
package smtp_provider

type Email struct {
	From    string
	To      []string
	CC      []string
	Subject string
	Body    string
}

type EmailProvider interface {
	SendEmail(email Email) error
}
