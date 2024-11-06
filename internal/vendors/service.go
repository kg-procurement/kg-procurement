//go:generate mockgen -typed -source=service.go -destination=service_mock.go -package=vendors
package vendors

import (
	"context"
	"fmt"
	"kg/procurement/cmd/config"
	"kg/procurement/cmd/utils"
	"kg/procurement/internal/common/database"
	"kg/procurement/internal/mailer"
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
	smtpProvider mailer.EmailProvider
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

func (v *VendorService) BlastEmail(ctx context.Context, vendorIDs []string, template emailTemplate) ([]string, error) {
	vendors, err := v.vendorDBAccessor.BulkGetByIDs(ctx, vendorIDs)
	if err != nil {
		return nil, err
	}

	errCh := make(chan error, len(vendors))
	defer close(errCh)

	// limit the number of concurrent workers to 20
	workerLimit := 20
	sem := make(chan struct{}, workerLimit)

	var wg sync.WaitGroup

	for _, vendor := range vendors {
		wg.Add(1)

		sem <- struct{}{}

		go func(vendor Vendor) {
			defer wg.Done()

			// replaces {{name}} keyword to vendor name
			bodyWithVendorName := strings.Replace(template.Body, "{{name}}", vendor.Name, -1)
			err := v.smtpProvider.SendEmail(mailer.Email{
				From:    v.cfg.SMTP.AuthEmail,
				To:      []string{vendor.Email},
				Subject: template.Subject,
				Body:    bodyWithVendorName,
			})
			errCh <- err

			<-sem
		}(vendor)
	}
	wg.Wait()

	var errList []string
	for i := 0; i < len(vendors); i++ {
		if err := <-errCh; err != nil {
			utils.Logger.Errorf("fail sending email %v", err)
			errList = append(errList, err.Error())
		}
	}

	if len(errList) > 0 {
		utils.Logger.Error("fail sending emails")
		return errList, fmt.Errorf("fail sending emails")
	}

	return nil, nil
}

func NewVendorService(
	cfg config.Application,
	conn database.DBConnector,
	clock clock.Clock,
	smtpProvider mailer.EmailProvider,
) *VendorService {
	return &VendorService{
		cfg:              cfg,
		vendorDBAccessor: newPostgresVendorAccessor(conn, clock),
		smtpProvider:     smtpProvider,
	}
}
