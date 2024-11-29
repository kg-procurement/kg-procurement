//go:generate mockgen -typed -source=mailer.go -destination=mailer_mock.go -package=mailer
package mailer

import (
	"kg/procurement/internal/common/database"
	"time"
)

type Email struct {
	From    string
	To      []string
	CC      []string
	Subject string
	Body    string
}

type EmailStatus struct {
	ID           string    `db:"id" json:"id"`
	EmailTo      string    `db:"email_to" json:"email_to"`
	Status       string    `db:"status" json:"status"`
	VendorID     string    `db:"vendor_id" json:"vendor_id"`
	DateSent     time.Time `db:"date_sent" json:"date_sent"`
	ModifiedDate time.Time `db:"modified_date" json:"modified_date"`
}

type EmailProvider interface {
	SendEmail(email Email) error
}

type GetAllEmailStatusSpec struct {
	ID           string    `db:"id" json:"id"`
	EmailTo      string    `db:"email_to" json:"email_to"`
	Status       string    `db:"status" json:"status"`
	ModifiedDate time.Time `db:"modified_date" json:"modified_date"`
	database.PaginationSpec
}

type AccessorGetEmailStatusPaginationData struct {
	EmailStatus []EmailStatus               `json:"email_status"`
	Metadata    database.PaginationMetadata `json:"metadata"`
}