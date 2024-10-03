//go:generate mockgen -typed -source=service.go -destination=service_mock.go -package=vendors
package vendors

import (
	"context"
	"kg/procurement/internal/common/database"

	"github.com/benbjohnson/clock"
)

type vendorDBAccessor interface {
	GetSomeStuff(ctx context.Context) ([]string, error)
	GetAll(ctx context.Context, spec GetAllVendorSpec) (*AccessorGetAllPaginationData, error)
	GetById(ctx context.Context, id string) (*Vendor, error)
	Put(ctx context.Context, spec Vendor) (*Vendor, error)
}

type VendorService struct {
	vendorDBAccessor
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

func (v *VendorService) Put(ctx context.Context, vendor Vendor) (*Vendor, error) {
	return v.vendorDBAccessor.Put(ctx, vendor)
}

func NewVendorService(
	conn database.DBConnector,
	clock clock.Clock,
) *VendorService {
	return &VendorService{
		vendorDBAccessor: newPostgresVendorAccessor(conn, clock),
	}
}
