//go:generate mockgen -typed -source=service.go -destination=service_mock.go -package=product
package product

import (
	"context"
	"kg/procurement/internal/common/database"
	"kg/procurement/cmd/utils"

	"github.com/benbjohnson/clock"
)

type productDBAccessor interface {
	getProductCategoryByID(ctx context.Context, pvID string) (*ProductCategory, error)
	GetProductVendorsByVendor(ctx context.Context, vendorID string, spec GetProductVendorByVendorSpec) (*AccessorGetProductVendorsPaginationData, error)
	getProductByID(ctx context.Context, productID string) (*Product, error)
	getPriceByPVID(ctx context.Context, pvID string) (*Price, error)
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
) (*GetProductVendorsResponse, error) {
	productVendors, err := p.productDBAccessor.GetProductVendorsByVendor(ctx, vendorID, spec)
	if err != nil {
		utils.Logger.Errorf(err.Error())
		return nil, err
	}
	return p.buildProductVendorsResponse(ctx, productVendors)
}

func (p *ProductService) GetProductVendors(
	ctx context.Context,
	spec GetProductVendorsSpec,
) (*GetProductVendorsResponse, error) {
	productVendors, err := p.productDBAccessor.GetAllProductVendors(ctx, spec)
	if err != nil {
		return nil, err
	}
	return p.buildProductVendorsResponse(ctx, productVendors)
}

func (p *ProductService) GetProductVendors(
	ctx context.Context,
	spec GetProductVendorsSpec,
) (*GetProductVendorsResponse, error) {
	productVendors, err := p.productDBAccessor.GetAllProductVendors(ctx, spec)
	if err != nil {
		utils.Logger.Errorf(err.Error())
		return nil, err
	}
	return p.buildProductVendorsResponse(ctx, productVendors)
}

func (p *ProductService) buildProductVendorsResponse(
	ctx context.Context,
	productVendors *AccessorGetProductVendorsPaginationData,
) (*GetProductVendorsResponse, error) {
	res := GetProductVendorsResponse{}
	for _, pv := range productVendors.ProductVendors {
		product, err := p.getProductByID(ctx, pv.ProductID)
		if err != nil {
			utils.Logger.Errorf(err.Error())
			return nil, err
		}

		category, err := p.getProductCategoryByID(ctx, product.ProductCategoryID)
		if err != nil {
			utils.Logger.Errorf(err.Error())
			return nil, err
		}

		price, err := p.getPriceByPVID(ctx, pv.ID)
		if err != nil {
			utils.Logger.Errorf(err.Error())
			return nil, err
		}

		pvr := ToProductVendorResponse(&pv, product, price, category)
    
		res.ProductVendors = append(res.ProductVendors, *pvr)
	}

	res.Metadata = productVendors.Metadata
	return &res, nil
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
		utils.Logger.Errorf(err.Error())
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
