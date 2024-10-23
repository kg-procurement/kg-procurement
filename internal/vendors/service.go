//go:generate mockgen -typed -source=service.go -destination=service_mock.go -package=vendors
package vendors

import (
	"context"
	"kg/procurement/cmd/config"
	"kg/procurement/internal/common/database"
	"kg/procurement/internal/smtp_provider"
	"log"
	"strings"
	"sync"

	"github.com/benbjohnson/clock"
)

type vendorDBAccessor interface {
	GetSomeStuff(ctx context.Context) ([]string, error)
	GetAll(ctx context.Context, spec GetAllVendorSpec) (*AccessorGetAllPaginationData, error)
	GetById(ctx context.Context, id string) (*Vendor, error)
	UpdateDetail(ctx context.Context, spec Vendor) (*Vendor, error)
	GetAllLocations(ctx context.Context) ([]string, error)
	BulkGetByIDs(_ context.Context, ids []string) ([]Vendor, error)
}

type VendorService struct {
	cfg config.Application
	vendorDBAccessor
	smtpProvider smtp_provider.EmailProvider
}

func (v *VendorService) GetSomeStuff(ctx context.Context) ([]string, error) {
	return v.vendorDBAccessor.GetSomeStuff(ctx)
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

func (v *VendorService) BlastEmail(ctx context.Context, vendorIDs []string, template emailTemplate) error {
	vendors, err := v.vendorDBAccessor.BulkGetByIDs(ctx, vendorIDs)
	if err != nil {
		return err
	}

	errCh := make(chan error, len(vendors))
	defer close(errCh)

	var wg sync.WaitGroup

	for _, vendor := range vendors {
		wg.Add(1)
		go func(vendor Vendor) {
			defer wg.Done()

			// replaces {{name}} keyword to vendor name
			bodyWithVendorName := strings.Replace(template.Body, "{{name}}", vendor.Name, -1)
			err := v.smtpProvider.SendEmail(smtp_provider.Email{
				From:    v.cfg.SMTP.AuthEmail,
				To:      []string{vendor.Email},
				Subject: template.Subject,
				Body:    bodyWithVendorName,
			})
			errCh <- err
		}(vendor)
	}
	wg.Wait()

	for i := 0; i < len(vendors); i++ {
		if err := <-errCh; err != nil {
			log.Printf("ERROR: fail sending email %v", err)
		}
	}
	return nil
}

func NewVendorService(
	cfg config.Application,
	conn database.DBConnector,
	clock clock.Clock,
	smtpProvider smtp_provider.EmailProvider,
) *VendorService {
	return &VendorService{
		cfg:              cfg,
		vendorDBAccessor: newPostgresVendorAccessor(conn, clock),
		smtpProvider:     smtpProvider,
	}
}
