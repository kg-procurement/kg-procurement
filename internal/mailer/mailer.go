//go:generate mockgen -typed -source=mailer.go -destination=mailer_mock.go -package=mailer
package mailer

import (
	"errors"
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
	ModifiedDate time.Time `db:"modified_date" json:"modified_date"`
}

type EmailProvider interface {
	SendEmail(email Email) error
}

type GetAllEmailStatusSpec struct {
	ID           string    `db:"id" json:"id"`
	EmailTo      string    `db:"email_to"`
	Status       string    `db:"status"`
	ModifiedDate time.Time `db:"modified_date"`
	database.PaginationSpec
}

type EmailStatusEnum int64

const (
	Success EmailStatusEnum = iota
	Failed
	InProgress
	Completed
)

func (s EmailStatusEnum) String() string {
	switch s {
	case Success:
		return "success"
	case Failed:
		return "failed"
	case InProgress:
		return "in_progress"
	case Completed:
		return "completed"
	}
	return "unknown"
}

func ParseEmailStatusEnum(status string) (EmailStatusEnum, error) {
	switch status {
	case "success":
		return Success, nil
	case "failed":
		return Failed, nil
	case "in_progress":
		return InProgress, nil
	case "completed":
		return Completed, nil
	default:
		return -1, errors.New("invalid email status")
	}
}
