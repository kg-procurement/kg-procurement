package product

import (
	"kg/procurement/internal/common/database"
	"kg/procurement/internal/vendors"
	"time"
)

type ProductID string

type Product struct {
	ID                ProductID `json:"id"`
	ProductCategoryID string    `json:"product_category_id"`
	UOMID             string    `json:"uom_id"`
	IncomeTaxID       string    `json:"income_tax_id"`
	ProductTypeID     string    `json:"product_type_id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	ModifiedDate      time.Time `json:"modified_date"`
	ModifiedBy        string    `json:"modified_by"`
}

type ProductCategory struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	ParentID       string    `json:"parent_id"`
	SpecialistBPID string    `json:"specialist_bpid"`
	ModifiedDate   time.Time `json:"modified_date"`
	ModifiedBy     string    `json:"modified_by"`
}

type ProductType struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Goods        bool      `json:"goods"`
	Asset        bool      `json:"asset"`
	Stockable    bool      `json:"stockable"`
	ModifiedDate time.Time `json:"modified_date"`
	ModifiedBy   string    `json:"modified_by"`
}

type UOM struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Dimension    string    `json:"dimension"`
	SAPCode      string    `json:"sap_code"`
	ModifiedDate time.Time `json:"modified_date"`
	ModifiedBy   string    `json:"modified_by"`
	StatusID     string    `json:"status_id"`
}

type ProductVendor struct {
	Product
	vendors.Vendor
}

type Price struct {
	ID                 string    `json:"id"`
	PurchasingOrgID    string    `json:"purchasing_org_id"`
	VendorID           string    `json:"vendor_id"`
	ProductVendorID    string    `json:"product_vendor_id"`
	QuantityMin        int       `json:"quantity_min"`
	QuantityMax        int       `json:"quantity_max"`
	QuantityUOMID      string    `json:"quantity_uom_id"`
	LeadTimeMin        int       `json:"lead_time_min"`
	LeadTimeMax        int       `json:"lead_time_max"`
	CurrencyID         string    `json:"currency_id"`
	Price              float64   `json:"price"`
	PriceQuantity      int       `json:"price_quantity"`
	PriceUOMID         string    `json:"price_uom_id"`
	ValidFrom          time.Time `json:"valid_from"`
	ValidTo            time.Time `json:"valid_to"`
	ValidPatternID     string    `json:"valid_pattern_id"`
	AreaGroupID        string    `json:"area_group_id"`
	ReferenceNumber    string    `json:"reference_number"`
	ReferenceDate      time.Time `json:"reference_date"`
	DocumentTypeID     string    `json:"document_type_id"`
	DocumentID         string    `json:"document_id"`
	ItemID             string    `json:"item_id"`
	TermOfPaymentID    string    `json:"term_of_payment_id"`
	InvocationOrder    int       `json:"invocation_order"`
	ModifiedDate       time.Time `json:"modified_date"`
	ModifiedBy         string    `json:"modified_by"`
}
type GetProductsByVendorSpec struct {
	Name     string `json:"name"`
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
    Name               string `json:"name"`
    PurchasingOrgID     string `json:"purchasing_org_id"`
    VendorID           string `json:"vendor_id"`
    ProductVendorID    string `json:"product_vendor_id"`
    QuantityMin        int    `json:"quantity_min"`
    QuantityMax        int    `json:"quantity_max"`
    QuantityUOMID      string `json:"quantity_uom_id"`
    LeadTimeMin        int    `json:"lead_time_min"`
    LeadTimeMax        int    `json:"lead_time_max"`
    CurrencyID         string `json:"currency_id"`
    Price              float64 `json:"price"`
    PriceQuantity      int    `json:"price_quantity"`
    PriceUOMID         string `json:"price_uom_id"`
    ValidFrom          time.Time `json:"valid_from"`
    ValidTo            time.Time `json:"valid_to"`
    ValidPatternID     string `json:"valid_pattern_id"`
    AreaGroupID        string `json:"area_group_id"`
    ReferenceNumber    string `json:"reference_number"`
    ReferenceDate      time.Time `json:"reference_date"`
    DocumentTypeID     string `json:"document_type_id"`
    DocumentID         string `json:"document_id"`
    ItemID             string `json:"item_id"`
    TermOfPaymentID    string `json:"term_of_payment_id"`
    InvocationOrder    int    `json:"invocation_order"`
}


