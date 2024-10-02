//go:generate mockgen -typed -source=service.go -destination=service_mock.go -package=vendors
package vendors

import (
	"context"
	"kg/procurement/internal/common/database"
	"strings"
)

type vendorDBAccessor interface {
	GetSomeStuff(ctx context.Context) ([]string, error)
	GetAll(ctx context.Context, spec database.PaginationSpec) (*AccessorGetAllPaginationData, error)
	GetByLocation(ctx context.Context, location string) ([]Vendor, error)
	GetByProductDescription(ctx context.Context, productDescription []string) ([]Vendor, error)
	GetById(ctx context.Context, id string) (*Vendor, error)
	Put(ctx context.Context, spec Vendor) (*Vendor, error)
}

type VendorService struct {
	vendorDBAccessor
}

func (v *VendorService) GetSomeStuff(ctx context.Context) ([]string, error) {
	return v.vendorDBAccessor.GetSomeStuff(ctx)
}

func (v *VendorService) GetAll(ctx context.Context, spec database.PaginationSpec) (*AccessorGetAllPaginationData, error) {

	accessorSpec := database.PaginationSpec{
		Limit: spec.Limit,
		Page:  spec.Page,
		Order: spec.Order,
	}

	return v.vendorDBAccessor.GetAll(ctx, accessorSpec)

}

func (v *VendorService) GetByLocation(ctx context.Context, location string) ([]Vendor, error) {
	return v.vendorDBAccessor.GetByLocation(ctx, location)
}

func (v *VendorService) GetByProduct(ctx context.Context, product string) ([]Vendor, error) {
	productDescription := strings.Fields(product)
	return v.vendorDBAccessor.GetByProductDescription(ctx, productDescription)
}

func (v *VendorService) Put(ctx context.Context, vendor Vendor) (*Vendor, error) {
	currentVendor, err := v.vendorDBAccessor.GetById(ctx, vendor.ID)
	if err != nil {
		return nil, err
	}

	newVendor := Vendor(*currentVendor)
	newVendor.Name = vendor.Name
	newVendor.Description = vendor.Description
	newVendor.BpID = vendor.BpID
	newVendor.BpName = vendor.BpName
	newVendor.Rating = vendor.Rating
	newVendor.AreaGroupID = vendor.AreaGroupID
	newVendor.AreaGroupName = vendor.AreaGroupName
	newVendor.SapCode = vendor.SapCode

	return v.vendorDBAccessor.Put(ctx, newVendor)
}

func NewVendorService(
	conn database.DBConnector,
) *VendorService {
	return &VendorService{
		vendorDBAccessor: newPostgresVendorAccessor(conn),
	}
}
