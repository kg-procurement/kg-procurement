//go:generate mockgen -typed -source=service.go -destination=service_mock.go -package=vendors
package vendors

import "context"

type vendorDBAccessor interface {
	GetSomeStuff(ctx context.Context) ([]string, error)
	GetAll(ctx context.Context) ([]Vendor, error)
	GetByLocation(ctx context.Context, location string) ([]Vendor, error)
	GetByProduct(ctx context.Context, product string) ([]Vendor, error)
}

type VendorService struct {
	vendorDBAccessor vendorDBAccessor
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
	return v.vendorDBAccessor.GetByProduct(ctx, product)
}

func NewVendorService(vendorDBAccessor vendorDBAccessor) *VendorService {
	return &VendorService{
		vendorDBAccessor: vendorDBAccessor,
	}
}
