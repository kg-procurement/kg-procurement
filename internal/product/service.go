//go:generate mockgen -typed -source=service.go -destination=service_mock.go -package=product
package product

import (
	"context"
	"kg/procurement/internal/common/database"
)

type productDBAccessor interface {
	GetProductsByVendor(ctx context.Context, vendorID string) ([]Product, error)
	UpdatePrice(ctx context.Context, price Price) (*Price, error)
	UpdateProduct(ctx context.Context, payload Product) (*Product, error)
}

type ProductService struct {
	productDBAccessor
}

func (p *ProductService) GetProductsByVendor(ctx context.Context, vendorID string) ([]Product, error) {
	return p.productDBAccessor.GetProductsByVendor(ctx, vendorID)
}

func (p *ProductService) UpdateProduct(ctx context.Context, payload Product) (*Product, error) {
	return p.productDBAccessor.UpdateProduct(ctx, payload)
}

func (p *ProductService) UpdatePrice(ctx context.Context, price Price) (*Price, error) {
	return p.productDBAccessor.UpdatePrice(ctx, price)
}

func NewProductService(
	conn database.DBConnector,
) *ProductService {
	return &ProductService{
		productDBAccessor: newPostgresProductAccessor(conn),
	}
}
