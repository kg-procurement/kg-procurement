//go:generate mockgen -typed -source=service.go -destination=service_mock.go -package=vendors
package vendors

import "context"

type vendorDBAccessor interface {
	GetSomeStuff(ctx context.Context) ([]string, error)
}

type VendorService struct {
	vendorDBAccessor vendorDBAccessor
}

func (v *VendorService) GetSomeStuff(ctx context.Context) ([]string, error) {
	return v.vendorDBAccessor.GetSomeStuff(ctx)
}

func NewVendorService(vendorDBAccessor vendorDBAccessor) *VendorService {
	return &VendorService{
		vendorDBAccessor: vendorDBAccessor,
	}
}
