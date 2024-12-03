//go:generate mockgen -typed -source=mailer.go -destination=mailer_mock.go -package=mailer
package mailer

import (
	"io"
	"kg/procurement/cmd/utils"
	"errors"
	"kg/procurement/internal/common/database"
	"mime/multipart"
	"time"
)

type Email struct {
	From        string
	To          []string
	CC          []string
	Subject     string
	Body        string
	Attachments []Attachment
}

type Attachment struct {
	Filename string
	Data     []byte
	MIMEType string
}

func BulkFromMultipartForm(files []*multipart.FileHeader) ([]Attachment, error) {
	var attachments []Attachment
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			utils.Logger.Error(err.Error())
			return nil, err
		}
		defer file.Close()

		fileData, err := io.ReadAll(file)
		if err != nil {
			utils.Logger.Error(err.Error())
			return nil, err
		}

		attachments = append(attachments, Attachment{
			Filename: fileHeader.Filename,
			Data:     fileData,
			MIMEType: fileHeader.Header.Get("Content-Type"),
		})
	}

	return attachments, nil
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
	EmailTo      string    `db:"email_to" json:"email_to"`
	Status       string    `db:"status" json:"status"`
	ModifiedDate time.Time `db:"modified_date" json:"modified_date"`
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
