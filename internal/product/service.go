//go:generate mockgen -typed -source=service.go -destination=service_mock.go -package=product
package product

import (
	"context"
	"errors"
	"kg/procurement/internal/common/database"
)

type productDBAccessor interface {
	GetProductsByVendor(ctx context.Context, vendorID string) ([]Product, error)
}

type ProductService struct {
	productDBAccessor
}

func (p *ProductService) GetProductsByVendor(_ context.Context, _ string) ([]Product, error) {
	return nil, errors.New("implement me")
}

func NewProductService(
	conn database.DBConnector,
) *ProductService {
	return &ProductService{
		productDBAccessor: newPostgresProductAccessor(conn),
	}
}
