//go:generate mockgen -typed -source=service.go -destination=service_mock.go -package=vendors
package vendors

import (
	"context"
	"encoding/json"
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
	CreateEvaluation(ctx context.Context, evaluation *VendorEvaluation) (*VendorEvaluation, error)
}

type emailStatusSvc interface {
	WriteEmailStatus(ctx context.Context, status mailer.EmailStatus) error
}

type VendorService struct {
	cfg config.Application
	vendorDBAccessor
	smtpProvider   mailer.EmailProvider
	emailStatusSvc emailStatusSvc
	redisClient    database.RedisClientInterface
}

func (v *VendorService) GetById(ctx context.Context, id string) (*Vendor, error) {
	cacheKey := fmt.Sprintf("vendor:%s", id)
	if cachedData, err := v.redisClient.Get(ctx, cacheKey); err == nil {
		var vendor Vendor
		if err := json.Unmarshal([]byte(cachedData), &vendor); err == nil {
			// cache hit
			return &vendor, nil
		}
	}

	// cache miss
	vendor, err := v.vendorDBAccessor.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if vendorJSON, err := json.Marshal(vendor); err == nil {
		_ = v.redisClient.Set(ctx, cacheKey, vendorJSON, 0)
	}

	return vendor, nil
}

func (v *VendorService) GetAll(ctx context.Context, spec GetAllVendorSpec) (*AccessorGetAllPaginationData, error) {
	return v.vendorDBAccessor.GetAll(ctx, spec)
}

func (v *VendorService) UpdateDetail(ctx context.Context, vendor Vendor) (*Vendor, error) {
	updatedVendor, err := v.vendorDBAccessor.UpdateDetail(ctx, vendor)
	if err != nil {
		return nil, err
	}

	cacheKey := fmt.Sprintf("vendor:%s", updatedVendor.ID)
	v.redisClient.Delete(ctx, cacheKey)

	return updatedVendor, nil
}

func (v *VendorService) GetLocations(ctx context.Context) ([]string, error) {
	return v.vendorDBAccessor.GetAllLocations(ctx)
}

func (v *VendorService) BlastEmail(ctx context.Context, vendorIDs []string, email mailer.Email) ([]string, error) {
	vendors, err := v.vendorDBAccessor.BulkGetByIDs(ctx, vendorIDs)
	if err != nil {
		return nil, err
	}

	v.applyDefaultEmailTemplate(&email)

	return v.executeBlastEmail(ctx, vendors, email)
}

func (v *VendorService) AutomatedEmailBlast(ctx context.Context, productName string) ([]string, error) {
	vendors, err := v.vendorDBAccessor.BulkGetByProductName(ctx, productName)

	if err != nil {
		return nil, err
	}

	email := &mailer.Email{}
	v.applyDefaultEmailTemplate(email)
	replacements := map[string]string{
		"{{product_name}}": productName,
	}

	email.Body = v.replacePlaceholder(email.Body, replacements)

	return v.executeBlastEmail(ctx, vendors, *email)
}

func (*VendorService) applyDefaultEmailTemplate(email *mailer.Email) {
	if email.Subject == "" {
		email.Subject = "Request for products"
	}
	if email.Body == "" {
		email.Body = "Kepada Yth {{name}},\n\nKami mengajukan permintaan untuk pengadaan produk {{product_name}} yang dibutuhkan oleh perusahaan kami. Mohon informasi mengenai ketersediaan, harga, dan waktu pengiriman untuk produk tersebut.\n\nTerima kasih atas perhatian dan kerjasamanya.\n\nHormat kami"
	}
}
    
func (v *VendorService) CreateEvaluation(ctx context.Context, evaluation *VendorEvaluation) (*VendorEvaluation, error) {
	id, _ := helper.GenerateRandomID()
	evaluation.ID = id

	evaluation.ModifiedDate = time.Now()

	return v.vendorDBAccessor.CreateEvaluation(ctx, evaluation)
}
    
func (v *VendorService) executeBlastEmail(ctx context.Context, vendors []Vendor, email mailer.Email) ([]string, error) {
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

			em := mailer.Email{
				From:        v.cfg.SMTP.AuthEmail,
				To:          []string{vendor.Email},
				CC:          email.CC,
				Subject:     email.Subject,
				Body:        v.replacePlaceholder(email.Body, replacements),
				Attachments: email.Attachments,
			}

			sendErr := v.smtpProvider.SendEmail(em)

			id, _ := helper.GenerateRandomID()
			emailStatus := mailer.EmailStatus{
				ID:           id,
				EmailTo:      vendor.Email,
				ModifiedDate: time.Now(),
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

func NewVendorService(
	cfg config.Application,
	conn database.DBConnector,
	clock clock.Clock,
	smtpProvider mailer.EmailProvider,
	emailStatusSvc emailStatusSvc,
	redisClient database.RedisClientInterface,
) *VendorService {
	return &VendorService{
		cfg:              cfg,
		vendorDBAccessor: newPostgresVendorAccessor(conn, clock),
		smtpProvider:     smtpProvider,
		emailStatusSvc:   emailStatusSvc,
		redisClient:      redisClient,
	}
}
