package mailer

import (
	"kg/procurement/internal/common/database"
)

type AccessorGetAllPaginationData struct {
	EmailStatus []EmailStatus               `json:"email_status"`
	Metadata    database.PaginationMetadata `json:"metadata"`
}
