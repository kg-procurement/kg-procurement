//go:generate mockgen -typed -source=service.go -destination=service_mock.go -package=product
package product

import (
	"context"
	u "kg/procurement/cmd/utils"
	"kg/procurement/internal/common/database"

	"github.com/benbjohnson/clock"
)

type productDBAccessor interface {
	GetProductsByVendor(ctx context.Context, vendorID string, spec GetProductsByVendorSpec) (*AccessorGetProductsByVendorPaginationData, error)
	UpdatePrice(ctx context.Context, price Price) (Price, error)
	UpdateProduct(ctx context.Context, payload Product) (Product, error)
}

type ProductService struct {
	productDBAccessor
}

func (p *ProductService) GetProductsByVendor(
	ctx context.Context,
	vendorID string,
	spec GetProductsByVendorSpec,
) (*AccessorGetProductsByVendorPaginationData, error) {
	u.GeneralLogger.Println("Starting getProductsByVendor process in service layer")

	result, err := p.productDBAccessor.GetProductsByVendor(ctx, vendorID, spec)
	if err != nil {
		u.ErrorLogger.Println(err.Error())
	}

	u.GeneralLogger.Println("ProductsByVendor successfully fetched from service layer")

	return result, err
}

func (p *ProductService) UpdateProduct(ctx context.Context, payload Product) (Product, error) {
	u.GeneralLogger.Println("Starting updateProduct process in service layer")

	result, err := p.productDBAccessor.UpdateProduct(ctx, payload)

	if err != nil {
		u.ErrorLogger.Println(err.Error())
	}

	u.ErrorLogger.Println("Product successfully udpated from service layer")

	return result, err
}

func (p *ProductService) UpdatePrice(ctx context.Context, price Price) (Price, error) {
	u.GeneralLogger.Println("Starting updateProductPrice process in service layer")

	result, err := p.productDBAccessor.UpdatePrice(ctx, price)

	if err != nil {
		u.ErrorLogger.Println(err.Error())
	}

	return result, err
}

func NewProductService(
	conn database.DBConnector,
	clock clock.Clock,
) *ProductService {
	return &ProductService{
		productDBAccessor: newPostgresProductAccessor(conn, clock),
	}
}
