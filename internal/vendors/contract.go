package vendors

type EmailBlastContract struct {
	VendorIDs     []string      `json:"vendor_ids" binding:"required"`
	EmailTemplate emailTemplate `json:"email_template" binding:"required"`
}

type emailTemplate struct {
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
}
