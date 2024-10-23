//go:generate mockgen -typed -source=mailer.go -destination=mailer_mock.go -package=mailer
package mailer

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
