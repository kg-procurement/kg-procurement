//go:generate mockgen -typed -source=service.go -destination=service_mock.go -package=vendors
package vendors

import (
	"context"
	"kg/procurement/internal/common/database"
)

type vendorDBAccessor interface {
	GetSomeStuff(ctx context.Context) ([]string, error)
	GetAll(ctx context.Context, spec database.PaginationSpec) (*AccessorGetAllPaginationData, error)
	GetByLocation(ctx context.Context, location string) ([]Vendor, error)
}

type VendorService struct {
	vendorDBAccessor
}

func (v *VendorService) GetSomeStuff(ctx context.Context) ([]string, error) {
	return v.vendorDBAccessor.GetSomeStuff(ctx)
}

func (v *VendorService) GetAll(ctx context.Context, spec database.PaginationSpec) (*ServiceGetAllPaginationData, error) {

	accessorSpec := database.PaginationSpec{
		Limit: spec.Limit,
		Page:  spec.Page,
		Order: spec.Order,
	}

	result, err := v.vendorDBAccessor.GetAll(ctx, accessorSpec)
	if err != nil {
		return nil, err
	}

	payload := ServiceGetAllPaginationData{
		Vendors:  result.Vendors,
		Metadata: result.Metadata,
	}

	return &payload, nil
}

func (v *VendorService) GetByLocation(ctx context.Context, location string) ([]Vendor, error) {
	return v.vendorDBAccessor.GetByLocation(ctx, location)
}

func NewVendorService(
	conn database.DBConnector,
) *VendorService {
	return &VendorService{
		vendorDBAccessor: newPostgresVendorAccessor(conn),
	}
}
