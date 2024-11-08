package product

import (
	"kg/procurement/internal/common/database"
	"time"
)

type ProductCategoryResponse struct {
	ID           string    `json:"id"`
	CategoryName string    `json:"category_name"`
	Description  string    `json:"description"`
	ModifiedDate time.Time `json:"modified_date"`
	ModifiedBy   string    `json:"modified_by"`
}

func newFromProductCategory(product_category *ProductCategory) *ProductCategoryResponse {
	return &ProductCategoryResponse{
		ID:                product_category.ID,
		CategoryName:      product_category.Name,
		Description:       product_category.Description,
		ModifiedDate:      product_category.ModifiedDate,
		ModifiedBy:        product_category.ModifiedBy,
	}
}
type ProductResponse struct {
	ID                ProductID                `json:"id"`
	ProductCategory   ProductCategoryResponse  `json:"product_category"`
	UOMID             string                   `json:"uom_id"`
	IncomeTaxID       string                   `json:"income_tax_id"`
	ProductTypeID     string                   `json:"product_type_id"`
	Name              string                   `json:"name"`
	Description       string                   `json:"description"`
	ModifiedDate      time.Time                `json:"modified_date"` // parse as time.dateTime
	ModifiedBy        string                   `json:"modified_by"`
}

func newFromProduct(product *Product, product_category *ProductCategory) *ProductResponse {
	productCategoryResponse := newFromProductCategory(product_category)
	return &ProductResponse{
		ID:                product.ID,
		ProductCategory:   *productCategoryResponse,
		UOMID:             product.UOMID,
		IncomeTaxID:       product.IncomeTaxID,
		ProductTypeID:     product.ProductTypeID,
		Name:              product.Name,
		Description:       product.Description,
		ModifiedDate:      product.ModifiedDate,
		ModifiedBy:        product.ModifiedBy,
	}
}

type ProductVendorResponse struct {
	ID                  string          `json:"id"`
	Product             ProductResponse `json:"product"`
	Code                string          `json:"code"`
	Name                string          `json:"name"`
	IncomeTaxID         string          `json:"income_tax_id"`
	IncomeTaxName       string          `json:"income_tax_name"`
	IncomeTaxPercentage string          `json:"income_tax_percentage"`
	Description         string          `json:"description"`
	UOMID               string          `json:"uom_id"`
	SAPCode             string          `json:"sap_code"`
	ModifiedDate        time.Time       `json:"modified_date"`
	ModifiedBy          string          `json:"modified_by"`
}

func ToProductVendorResponse(pv *ProductVendor, p *Product, pc *ProductCategory) *ProductVendorResponse {
	productResponse := newFromProduct(p, pc)
	return &ProductVendorResponse{
		ID:                  pv.ID,
		Product:             *productResponse,
		Code:                pv.Code,
		Name:                pv.Name,
		IncomeTaxID:         pv.IncomeTaxID,
		IncomeTaxName:       pv.IncomeTaxName,
		IncomeTaxPercentage: pv.IncomeTaxPercentage,
		Description:         pv.Description,
		UOMID:               pv.UOMID,
		SAPCode:             pv.SAPCode,
		ModifiedBy:          pv.ModifiedBy,
		ModifiedDate:        pv.ModifiedDate,
	}
}

type GetProductVendorsByVendorResponse struct {
	ProductVendors []ProductVendorResponse     `json:"product_vendors"`
	Metadata       database.PaginationMetadata `json:"metadata"`
}
