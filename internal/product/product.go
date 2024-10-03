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

type ProductVendor struct {
	Product
	vendors.Vendor
}

type GetProductsByVendorSpec struct {
	Name     string `json:"name"`

	database.PaginationSpec
}
