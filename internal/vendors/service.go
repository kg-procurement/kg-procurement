//go:generate mockgen -typed -source=service.go -destination=service_mock.go -package=vendors
package vendors

import (
	"context"
	"kg/procurement/internal/common/database"
	"strings"
)

type vendorDBAccessor interface {
	GetSomeStuff(ctx context.Context) ([]string, error)
	GetAll(ctx context.Context, spec GetAllVendorSpec) (*AccessorGetAllPaginationData, error)
	GetByLocation(ctx context.Context, location string) ([]Vendor, error)
	GetByProductDescription(ctx context.Context, productDescription []string) ([]Vendor, error)
}

type VendorService struct {
	vendorDBAccessor
}

func (v *VendorService) GetSomeStuff(ctx context.Context) ([]string, error) {
	return v.vendorDBAccessor.GetSomeStuff(ctx)
}

func (v *VendorService) GetAll(ctx context.Context, spec GetAllVendorSpec) (*AccessorGetAllPaginationData, error) {
	return v.vendorDBAccessor.GetAll(ctx, spec)

}

func (v *VendorService) GetByLocation(ctx context.Context, location string) ([]Vendor, error) {
	return v.vendorDBAccessor.GetByLocation(ctx, location)
}

func (v *VendorService) GetByProduct(ctx context.Context, product string) ([]Vendor, error) {
	productDescription := strings.Fields(product)
	return v.vendorDBAccessor.GetByProductDescription(ctx, productDescription)
}

func NewVendorService(
	conn database.DBConnector,
) *VendorService {
	return &VendorService{
		vendorDBAccessor: newPostgresVendorAccessor(conn),
	}
}
