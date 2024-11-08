//go:generate mockgen -typed -source=service.go -destination=service_mock.go -package=product
package product

import (
	"context"
	"kg/procurement/internal/common/database"

	"github.com/benbjohnson/clock"
)

type productDBAccessor interface {
	getProductCategoryByID(ctx context.Context, pvID string) (*ProductCategory, error)
	GetProductVendorsByVendor(ctx context.Context, vendorID string, spec GetProductVendorByVendorSpec) (*AccessorGetProductVendorsByVendorPaginationData, error)
	getProductByID(ctx context.Context, productID string) (*Product, error)
	GetAllProductVendors(ctx context.Context, spec GetProductVendorsSpec) (*AccessorGetProductVendorsPaginationData, error)
	UpdatePrice(ctx context.Context, price Price) (Price, error)
	UpdateProduct(ctx context.Context, payload Product) (Product, error)
}

type ProductService struct {
	productDBAccessor
}

func (p *ProductService) GetProductVendorsByVendor(
	ctx context.Context,
	vendorID string,
	spec GetProductVendorByVendorSpec,
) (*GetProductVendorsByVendorResponse, error) {
	res := GetProductVendorsByVendorResponse{}
	productVendors, err := p.productDBAccessor.GetProductVendorsByVendor(ctx, vendorID, spec)
	if err != nil {
		return nil, err
	}

	for _, pv := range productVendors.ProductVendors {
		product, err := p.getProductByID(ctx, pv.ProductID)
		if err != nil {
			return nil, err
		}

		category, err := p.getProductCategoryByID(ctx, product.ProductCategoryID)
		if err != nil {
			return nil, err
		}

		pvr := ToProductVendorResponse(&pv, product, category)
		res.ProductVendors = append(res.ProductVendors, *pvr)
	}

	res.Metadata = productVendors.Metadata
	return &res, nil
}

func (p *ProductService) GetProductVendors(
	ctx context.Context,
	spec GetProductVendorsSpec,
) (*AccessorGetProductVendorsPaginationData, error) {
	return p.productDBAccessor.GetAllProductVendors(ctx, spec)
}

func (p *ProductService) UpdateProduct(ctx context.Context, payload Product) (Product, error) {
	return p.productDBAccessor.UpdateProduct(ctx, payload)
}

func (p *ProductService) UpdatePrice(ctx context.Context, price Price) (Price, error) {
	return p.productDBAccessor.UpdatePrice(ctx, price)
}

func (p *ProductService) getProductByID(ctx context.Context, productID string) (*Product, error) {
	product, err := p.productDBAccessor.getProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func NewProductService(
	conn database.DBConnector,
	clock clock.Clock,
) *ProductService {
	return &ProductService{
		productDBAccessor: newPostgresProductAccessor(conn, clock),
	}
}
