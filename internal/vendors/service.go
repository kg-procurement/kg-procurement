//go:generate mockgen -typed -source=service.go -destination=service_mock.go -package=vendors
package vendors

import (
	"context"
	"kg/procurement/internal/common/database"
)

type vendorDBAccessor interface {
	GetSomeStuff(ctx context.Context) ([]string, error)
	GetAll(ctx context.Context, spec database.GetAllPaginationSpec) (*AccessorGetAllPaginationData, error)
	GetByLocation(ctx context.Context, location string) ([]Vendor, error)
}

type VendorService struct {
	vendorDBAccessor
}

func (v *VendorService) GetSomeStuff(ctx context.Context) ([]string, error) {
	return v.vendorDBAccessor.GetSomeStuff(ctx)
}

func (v *VendorService) GetAll(ctx context.Context, spec ServiceGetAllPaginationSpec) (*ServiceGetAllPaginationData, error) {

	limit := 10
	if spec.Limit > 0 {
		limit = spec.Limit
	}

	offset := spec.Limit * (spec.Page - 1)

	accessorSpec := AccessorGetAllPaginationSpec{
		Limit:  limit,
		Offset: offset,
		Order:  spec.Order,
	}

	result, err := v.vendorDBAccessor.GetAll(ctx, accessorSpec)
	if err != nil {
		return nil, err
	}

	var previousPage *int = nil
	if spec.Page > 1 {
		previousPage = new(int)
		*previousPage = spec.Page - 1
	}

	payload := ServiceGetAllPaginationData{
		Vendors:      result.Vendors,
		TotalEntries: result.TotalEntries,
		CurrentPage:  spec.Page,
		PreviousPage: previousPage,
		NextPage:     spec.Page + 1,
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
