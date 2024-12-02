//go:generate mockgen -typed -source=service.go -destination=service_mock.go -package=vendors
package vendors

import (
	"context"
	"fmt"
	"kg/procurement/cmd/config"
	"kg/procurement/cmd/utils"
	"kg/procurement/internal/common/database"
	"kg/procurement/internal/common/helper"
	"kg/procurement/internal/mailer"
	"strings"
	"sync"
	"time"

	"github.com/benbjohnson/clock"
)

type vendorDBAccessor interface {
	GetSomeStuff(ctx context.Context) ([]string, error)
	GetAll(ctx context.Context, spec GetAllVendorSpec) (*AccessorGetAllPaginationData, error)
	GetById(ctx context.Context, id string) (*Vendor, error)
	UpdateDetail(ctx context.Context, spec Vendor) (*Vendor, error)
	GetAllLocations(ctx context.Context) ([]string, error)
	BulkGetByIDs(_ context.Context, ids []string) ([]Vendor, error)
	BulkGetByProductName(_ context.Context, productName string) ([]Vendor, error)
}

type emailStatusSvc interface {
	WriteEmailStatus(ctx context.Context, status mailer.EmailStatus) error
	GetAllEmailStatus(ctx context.Context, spec mailer.GetAllEmailStatusSpec) (*mailer.AccessorGetEmailStatusPaginationData, error)
}

type VendorService struct {
	cfg config.Application
	vendorDBAccessor
	smtpProvider   mailer.EmailProvider
	emailStatusSvc emailStatusSvc
}

func (v *VendorService) GetById(ctx context.Context, id string) (*Vendor, error) {
	return v.vendorDBAccessor.GetById(ctx, id)
}

func (v *VendorService) GetAll(ctx context.Context, spec GetAllVendorSpec) (*AccessorGetAllPaginationData, error) {
	return v.vendorDBAccessor.GetAll(ctx, spec)
}

func (v *VendorService) UpdateDetail(ctx context.Context, vendor Vendor) (*Vendor, error) {
	return v.vendorDBAccessor.UpdateDetail(ctx, vendor)
}

func (v *VendorService) GetLocations(ctx context.Context) ([]string, error) {
	return v.vendorDBAccessor.GetAllLocations(ctx)
}

func (v *VendorService) BlastEmail(ctx context.Context, vendorIDs []string, template emailTemplate) ([]string, error) {
	vendors, err := v.vendorDBAccessor.BulkGetByIDs(ctx, vendorIDs)
	if err != nil {
		return nil, err
	}

	v.applyDefaultEmailTemplate(&template)

	return v.executeBlastEmail(ctx, vendors, template)
}

func (v *VendorService) AutomatedEmailBlast(ctx context.Context, productName string) ([]string, error) {
	vendors, err := v.vendorDBAccessor.BulkGetByProductName(ctx, productName)

	if err != nil {
		return nil, err
	}

	template := &emailTemplate{}
	v.applyDefaultEmailTemplate(template)
	replacements := map[string]string{
		"{{product_name}}": productName,
	}

	template.Body = v.replacePlaceholder(template.Body, replacements)

	return v.executeBlastEmail(ctx, vendors, *template)
}

func (*VendorService) applyDefaultEmailTemplate(template *emailTemplate) {
	if template.Subject == "" {
		template.Subject = "Request for products"
	}
	if template.Body == "" {
		template.Body = "Kepada Yth {{name}},\n\nKami mengajukan permintaan untuk pengadaan produk {{product_name}} yang dibutuhkan oleh perusahaan kami. Mohon informasi mengenai ketersediaan, harga, dan waktu pengiriman untuk produk tersebut.\n\nTerima kasih atas perhatian dan kerjasamanya.\n\nHormat kami"
	}
}

func (v *VendorService) executeBlastEmail(ctx context.Context, vendors []Vendor, template emailTemplate) ([]string, error) {
	errCh := make(chan error, len(vendors))
	statusCh := make(chan mailer.EmailStatus, len(vendors))

	// Limit the number of concurrent workers to 20
	workerLimit := 20
	sem := make(chan struct{}, workerLimit)

	var wg sync.WaitGroup

	for _, vendor := range vendors {
		wg.Add(1)

		sem <- struct{}{}

		go func(vendor Vendor) {
			defer wg.Done()
			defer func() { <-sem }() // release the semaphore slot

			// replaces {{name}} keyword to vendor name
			replacements := map[string]string{
				"{{name}}": vendor.Name,
			}

			templateBody := v.replacePlaceholder(template.Body, replacements)

			email := mailer.Email{
				From:    v.cfg.SMTP.AuthEmail,
				To:      []string{vendor.Email},
				Subject: template.Subject,
				Body:    templateBody,
			}

			sendErr := v.smtpProvider.SendEmail(email)

			id, _ := helper.GenerateRandomID()
			dateSent := time.Now()
			emailStatus := mailer.EmailStatus{
				ID:           id,
				EmailTo:      vendor.Email,
				VendorID:     vendor.ID,
				DateSent:     dateSent,
				ModifiedDate: dateSent,
			}

			if sendErr != nil {
				emailStatus.Status = mailer.Failed.String()
				errCh <- sendErr
			} else {
				emailStatus.Status = mailer.Success.String()
				errCh <- nil
			}

			statusCh <- emailStatus
		}(vendor)
	}

	wg.Wait()
	close(errCh)
	close(statusCh)

	var errList []string
	for err := range errCh {
		if err != nil {
			utils.Logger.Errorf("failed to send email: %v", err)
			errList = append(errList, err.Error())
		}
	}

	// write email statuses to the database so we can track the status
	for status := range statusCh {
		writeErr := v.emailStatusSvc.WriteEmailStatus(ctx, status)
		if writeErr != nil {
			utils.Logger.Errorf("failed to write email status: %v", writeErr)
		}
	}

	if len(errList) > 0 {
		utils.Logger.Error("failed to send some emails")
		return errList, fmt.Errorf("failed to send some emails")
	}

	return nil, nil
}

func (v *VendorService) replacePlaceholder(template string, replacements map[string]string) string {
	for placeholder, value := range replacements {
		template = strings.ReplaceAll(template, placeholder, value)
	}
	return template
}

func (v *VendorService) GetPopulatedEmailStatus(
	ctx context.Context,
	spec mailer.GetAllEmailStatusSpec,
) (*mailer.GetAllEmailStatusResponse, error) {
	emailStatus, err := v.emailStatusSvc.GetAllEmailStatus(ctx, spec)
	if err != nil {
		utils.Logger.Errorf("Error fetching email statuses: %v", err)
		return nil, err
	}
	res := mailer.GetAllEmailStatusResponse{}
	var vendorIDs []string
	for _, es := range emailStatus.EmailStatus {
		vendorIDs = append(vendorIDs, es.VendorID)
	}

	vendors, err := v.BulkGetByIDs(ctx, vendorIDs)
	if err != nil {
		utils.Logger.Errorf("Error fetching vendors: %v", err)
		return nil, err
	}

	vendorMap := make(map[string]Vendor)
	for _, vendor := range vendors {
		vendorMap[vendor.ID] = vendor
	}

	for _, es := range emailStatus.EmailStatus {
		vendorName := "Unknown Vendor" 
		if vendor, exists := vendorMap[es.VendorID]; exists {
			vendorName = vendor.Name
		} else {
			utils.Logger.Infof("Vendor not found for VendorID: %s", es.VendorID)
		}

		emailStatusResponse := mailer.EmailStatusResponse{
			EmailStatus: mailer.EmailStatus{
				ID:           es.ID,
				EmailTo:      es.EmailTo,
				Status:       es.Status,
				VendorID:     es.VendorID,
				DateSent:     es.DateSent,
				ModifiedDate: es.ModifiedDate,
			},
			VendorName: vendorName,
		}

		res.EmailStatus = append(res.EmailStatus, emailStatusResponse)
	}

	res.Metadata = emailStatus.Metadata
	return &res, nil
}

func NewVendorService(
	cfg config.Application,
	conn database.DBConnector,
	clock clock.Clock,
	smtpProvider mailer.EmailProvider,
	emailStatusSvc emailStatusSvc,
) *VendorService {
	return &VendorService{
		cfg:              cfg,
		vendorDBAccessor: newPostgresVendorAccessor(conn, clock),
		smtpProvider:     smtpProvider,
		emailStatusSvc:   emailStatusSvc,
	}
}
