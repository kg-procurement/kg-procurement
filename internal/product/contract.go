package product

import (
	"kg/procurement/internal/common/database"
	"time"
)

type ProductResponse struct {
	ID                ProductID `json:"id"`
	ProductCategoryID string    `json:"product_category_id"`
	UOMID             string    `json:"uom_id"`
	IncomeTaxID       string    `json:"income_tax_id"`
	ProductTypeID     string    `json:"product_type_id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	ModifiedDate      time.Time `json:"modified_date"` // parse as time.dateTime
	ModifiedBy        string    `json:"modified_by"`
}

func newFromProduct(product *Product) *ProductResponse {
	return &ProductResponse{
		ID:                product.ID,
		ProductCategoryID: product.ProductCategoryID,
		UOMID:             product.UOMID,
		IncomeTaxID:       product.IncomeTaxID,
		ProductTypeID:     product.ProductTypeID,
		Name:              product.Name,
		Description:       product.Description,
		ModifiedDate:      product.ModifiedDate,
		ModifiedBy:        product.ModifiedBy,
	}
}

type PriceResponse struct {
	ID                string    `db:"id" json:"id"`
	PurchasingOrgID   string    `db:"purchasing_org_id" json:"purchasing_org_id"`
	PurchasingOrgName string    `db:"purchasing_org_name" json:"purchasing_org_name"`
	VendorID          string    `db:"vendor_id" json:"vendor_id"`
	ProductVendorID   string    `db:"product_vendor_id" json:"product_vendor_id"`
	QuantityMin       int       `db:"quantity_min" json:"quantity_min"`
	QuantityMax       int       `db:"quantity_max" json:"quantity_max"`
	QuantityUOMID     string    `db:"quantity_uom_id" json:"quantity_uom_id"`
	LeadTimeMin       int       `db:"lead_time_min" json:"lead_time_min"`
	LeadTimeMax       int       `db:"lead_time_max" json:"lead_time_max"`
	CurrencyID        string    `db:"currency_id" json:"currency_id"`
	CurrencyName      string    `db:"currency_name" json:"currency_name"`
	CurrencyCode      string    `db:"currency_code" json:"currency_code"`
	Price             float64   `db:"price" json:"price"`
	PriceQuantity     int       `db:"price_quantity" json:"price_quantity"`
	PriceUOMID        string    `db:"price_uom_id" json:"price_uom_id"`
	ValidFrom         time.Time `db:"valid_from" json:"valid_from"`
	ValidTo           time.Time `db:"valid_to" json:"valid_to"`
	ValidPatternID    string    `db:"valid_pattern_id" json:"valid_pattern_id"`
	ValidPatternName  string    `db:"valid_pattern_name" json:"valid_pattern_name"`
	AreaGroupID       string    `db:"area_group_id" json:"area_group_id"`
	AreaGroupName     string    `db:"area_group_name" json:"area_group_name"`
	ReferenceNumber   string    `db:"reference_number" json:"reference_number"`
	ReferenceDate     time.Time `db:"reference_date" json:"reference_date"`
	DocumentTypeID    string    `db:"document_type_id" json:"document_type_id"`
	DocumentTypeName  string    `db:"document_type_name" json:"document_type_name"`
	DocumentID        string    `db:"document_id" json:"document_id"`
	ItemID            string    `db:"item_id" json:"item_id"`
	TermOfPaymentID   string    `db:"term_of_payment_id" json:"term_of_payment_id"`
	TermOfPaymentDays int       `db:"term_of_payment_days" json:"term_of_payment_days"`
	TermOfPaymentText string    `db:"term_of_payment_text" json:"term_of_payment_text"`
	InvocationOrder   int       `db:"invocation_order" json:"invocation_order"`
	ModifiedDate      time.Time `db:"modified_date" json:"modified_date"`
	ModifiedBy        string    `db:"modified_by" json:"modified_by"`
}

func newFromPrice(price *Price) *PriceResponse {
	return &PriceResponse{
		ID:           price.ID,
		Price:        price.Price,
		ModifiedDate: price.ModifiedDate,
		ModifiedBy:   price.ModifiedBy,
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

func ToProductVendorResponse(pv *ProductVendor, p *Product, pr *Price) *ProductVendorResponse {
	productResponse := newFromProduct(p)
	priceResponse := newFromPrice(pr)
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

type GetProductVendorsByVendorResponse struct {
	ProductVendors []ProductVendorResponse     `json:"product_vendors"`
	Metadata       database.PaginationMetadata `json:"metadata"`
}
