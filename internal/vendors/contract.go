package vendors

import (
	"mime/multipart"
)

type EmailBlastContract struct {
	VendorIDs   []string                `form:"vendor_ids" binding:"required"`
	Subject     string                  `form:"subject"`
	Body        string                  `form:"body"`
	Attachments []*multipart.FileHeader `form:"attachments"`
}
