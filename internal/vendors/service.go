//go:generate mockgen -typed -source=service.go -destination=service_mock.go -package=vendors
package vendors

import (
	"context"
	"kg/procurement/internal/common/database"
	"strings"
)

type vendorDBAccessor interface {
	GetSomeStuff(ctx context.Context) ([]string, error)
	GetAll(ctx context.Context) ([]Vendor, error)
	GetByLocation(ctx context.Context, location string) ([]Vendor, error)
	GetByProductWords(ctx context.Context, productWords []string) ([]Vendor, error)
}

type VendorService struct {
	vendorDBAccessor
}

func (v *VendorService) GetSomeStuff(ctx context.Context) ([]string, error) {
	return v.vendorDBAccessor.GetSomeStuff(ctx)
}

func (v *VendorService) GetAll(ctx context.Context) ([]Vendor, error) {
	return v.vendorDBAccessor.GetAll(ctx)
}

func (v *VendorService) GetByLocation(ctx context.Context, location string) ([]Vendor, error) {
	return v.vendorDBAccessor.GetByLocation(ctx, location)
}

func (v *VendorService) GetByProduct(ctx context.Context, product string) ([]Vendor, error) {
	productWords := strings.Fields(product)
	return v.vendorDBAccessor.GetByProductWords(ctx, productWords)
}

func NewVendorService(
	conn database.DBConnector,
) *VendorService {
	return &VendorService{
		vendorDBAccessor: newPostgresVendorAccessor(conn),
	}
}
