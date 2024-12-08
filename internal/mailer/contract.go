package mailer

import (
	"kg/procurement/internal/common/database"
)

type GetAllEmailStatusResponse struct {
	EmailStatus []EmailStatusResponse       `json:"email_status"`
	Metadata    database.PaginationMetadata `json:"metadata"`
}

type EmailStatusResponse struct {
	EmailStatus
	VendorName string `json:"vendor_name"`
	VendorRating int `json:"vendor_rating"`
}
