package product

import (
	"kg/procurement/internal/common/database"
	"time"
)

type GetProductVendorsResponse struct {
	ProductVendors []ProductVendorResponse     `json:"product_vendors"`
	Metadata       database.PaginationMetadata `json:"metadata"`
}

type ProductCategoryResponse struct {
	ID           string    `json:"id"`
	CategoryName string    `json:"category_name"`
	Description  string    `json:"description"`
	ModifiedDate time.Time `json:"modified_date"`
	ModifiedBy   string    `json:"modified_by"`
}

func newFromProductCategory(productCategory *ProductCategory) *ProductCategoryResponse {
	return &ProductCategoryResponse{
		ID:           productCategory.ID,
		CategoryName: productCategory.Name,
		Description:  productCategory.Description,
		ModifiedDate: productCategory.ModifiedDate,
		ModifiedBy:   productCategory.ModifiedBy,
	}
}

type UOMResponse struct {
	ID      string `json:"id"`
	UOMName string `json:"uom_name"`
}

func newFromUOM(uom *UOM) *UOMResponse {
	return &UOMResponse{
		ID:           uom.ID,
		UOMName:      uom.Name,
	}
}

type ProductResponse struct {
	ID              ProductID               `json:"id"`
	ProductCategory ProductCategoryResponse `json:"product_category"`
	UOMID           string                  `json:"uom_id"`
	IncomeTaxID     string                  `json:"income_tax_id"`
	ProductTypeID   string                  `json:"product_type_id"`
	Name            string                  `json:"name"`
	Description     string                  `json:"description"`
	ModifiedDate    time.Time               `json:"modified_date"` // parse as time.dateTime
	ModifiedBy      string                  `json:"modified_by"`
}

func newProductResponseFromProduct(product *Product, productCategory *ProductCategory) *ProductResponse {
	productCategoryResponse := newFromProductCategory(productCategory)
	return &ProductResponse{
		ID:              product.ID,
		ProductCategory: *productCategoryResponse,
		UOMID:           product.UOMID,
		IncomeTaxID:     product.IncomeTaxID,
		ProductTypeID:   product.ProductTypeID,
		Name:            product.Name,
		Description:     product.Description,
		ModifiedDate:    product.ModifiedDate,
		ModifiedBy:      product.ModifiedBy,
	}
}

type PriceResponse struct {
	ID            string         `json:"id"`
	Price         float64        `json:"price"`
	CurrencyCode  string         `json:"currency_code"`
	PriceQuantity int         `json:"price_quantity"`
	VendorID      string         `json:"vendor_id"`
	UOM           UOMResponse    `json:"uom"`
	ModifiedDate  time.Time      `json:"modified_date"`
	ModifiedBy    string         `json:"modified_by"`
}

func newPriceResponseFromPrice(price *Price, uom *UOM) *PriceResponse {
	UOMResponse := newFromUOM(uom)
	return &PriceResponse{
		ID:            price.ID,
		Price:         price.Price,
		CurrencyCode:  price.CurrencyCode,
		PriceQuantity: price.PriceQuantity,
		VendorID:      price.VendorID,
		UOM:           *UOMResponse,
		ModifiedDate:  price.ModifiedDate,
		ModifiedBy:    price.ModifiedBy,
	}
}

type ProductVendorResponse struct {
	ID                  string          `json:"id"`
	Product             ProductResponse `json:"product"`
	Price               PriceResponse   `json:"price"`
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

func ToProductVendorResponse(pv *ProductVendor, p *Product, pr *Price, pc *ProductCategory, uom *UOM) *ProductVendorResponse {
	productResponse := newProductResponseFromProduct(p, pc)
	priceResponse := newPriceResponseFromPrice(pr, uom)
	return &ProductVendorResponse{
		ID:                  pv.ID,
		Product:             *productResponse,
		Price:               *priceResponse,
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
