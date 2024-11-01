package product

import (
	"kg/procurement/internal/common/database"
	"time"
)

type ProductID string

type Product struct {
	ID                ProductID `db:"id" json:"id"`
	ProductCategoryID string    `db:"product_category_id" json:"product_category_id"`
	UOMID             string    `db:"uom_id" json:"uom_id"`
	IncomeTaxID       string    `db:"income_tax_id" json:"income_tax_id"`
	ProductTypeID     string    `db:"product_type_id" json:"product_type_id"`
	Name              string    `db:"name" json:"name"`
	Description       string    `db:"description" json:"description"`
	ModifiedDate      time.Time `db:"modified_date" json:"modified_date"` // parse as time.dateTime
	ModifiedBy        string    `db:"modified_by" json:"modified_by"`
}

type ProductCategory struct {
	ID             string    `db:"id" json:"id"`
	Name           string    `db:"name" json:"name"`
	Code           string    `db:"code" json:"code"`
	Description    string    `db:"description" json:"description"`
	ParentID       string    `db:"parent_id" json:"parent_id"`
	SpecialistBPID string    `db:"specialist_bpid" json:"specialist_bpid"`
	ModifiedDate   time.Time `db:"modified_date" json:"modified_date"` // Will be parsed as time.Time later
	ModifiedBy     string    `db:"modified_by" json:"modified_by"`
}

type ProductType struct {
	ID           string    `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	Description  string    `db:"description" json:"description"`
	Goods        bool      `db:"goods" json:"goods"`
	Asset        bool      `db:"asset" json:"asset"`
	Stock        bool      `db:"stock" json:"stock"`
	ModifiedDate time.Time `db:"modified_date" json:"modified_date"`
	ModifiedBy   string    `db:"modified_by" json:"modified_by"`
}

type UOM struct {
	ID           string    `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	Description  string    `db:"description" json:"description"`
	Dimension    string    `db:"dimension" json:"dimension"`
	SAPCode      string    `db:"sap_code" json:"sap_code"`
	ModifiedDate time.Time `db:"modified_date" json:"modified_date"`
	ModifiedBy   string    `db:"modified_by" json:"modified_by"`
	StatusID     string    `db:"status_id" json:"status_id"`
}

type ProductVendor struct {
	ID                  string    `db:"id" json:"id"`
	ProductID           string    `db:"product_id" json:"product_id"`
	Code                string    `db:"code" json:"code"`
	Name                string    `db:"name" json:"name"`
	IncomeTaxID         string    `db:"income_tax_id" json:"income_tax_id"`
	IncomeTaxName       string    `db:"income_tax_name" json:"income_tax_name"`
	IncomeTaxPercentage string    `db:"income_tax_percentage" json:"income_tax_percentage"`
	Description         string    `db:"description" json:"description"`
	UOMID               string    `db:"uom_id" json:"uom_id"`
	SAPCode             string    `db:"sap_code" json:"sap_code"`
	ModifiedDate        time.Time `db:"modified_date" json:"modified_date"`
	ModifiedBy          string    `db:"modified_by" json:"modified_by"`
}

type Price struct {
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

type GetProductVendorByVendorSpec struct {
	Name string `json:"name"`
	database.PaginationSpec
}

type PutProductSpec struct {
	ProductCategoryID string `json:"product_category_id"`
	UOMID             string `json:"uom_id"`
	IncomeTaxID       string `json:"income_tax_id"`
	ProductTypeID     string `json:"product_type_id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
}

type PutPriceSpec struct {
	Name            string    `json:"name"`
	PurchasingOrgID string    `json:"purchasing_org_id"`
	VendorID        string    `json:"vendor_id"`
	ProductVendorID string    `json:"product_vendor_id"`
	QuantityMin     int       `json:"quantity_min"`
	QuantityMax     int       `json:"quantity_max"`
	QuantityUOMID   string    `json:"quantity_uom_id"`
	LeadTimeMin     int       `json:"lead_time_min"`
	LeadTimeMax     int       `json:"lead_time_max"`
	CurrencyID      string    `json:"currency_id"`
	Price           float64   `json:"price"`
	PriceQuantity   int       `json:"price_quantity"`
	PriceUOMID      string    `json:"price_uom_id"`
	ValidFrom       time.Time `json:"valid_from"`
	ValidTo         time.Time `json:"valid_to"`
	ValidPatternID  string    `json:"valid_pattern_id"`
	AreaGroupID     string    `json:"area_group_id"`
	ReferenceNumber string    `json:"reference_number"`
	ReferenceDate   time.Time `json:"reference_date"`
	DocumentTypeID  string    `json:"document_type_id"`
	DocumentID      string    `json:"document_id"`
	ItemID          string    `json:"item_id"`
	TermOfPaymentID string    `json:"term_of_payment_id"`
	InvocationOrder int       `json:"invocation_order"`
}
