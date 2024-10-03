//go:generate mockgen -typed -source=service.go -destination=service_mock.go -package=product
package product

import (
	"context"
	"kg/procurement/internal/common/database"
)

type productDBAccessor interface {
	GetProductsByVendor(ctx context.Context, spec GetProductsByVendorSpec) ([]Product, error)
}

type ProductService struct {
	productDBAccessor
}

func (p *ProductService) GetProductsByVendor(ctx context.Context, spec GetProductsByVendorSpec) ([]Product, error) {
	return p.productDBAccessor.GetProductsByVendor(ctx, spec)
}

func NewProductService(
	conn database.DBConnector,
) *ProductService {
	return &ProductService{
		productDBAccessor: newPostgresProductAccessor(conn),
	}
}
