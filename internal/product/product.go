package product

import (
	"kg/procurement/internal/common/database"
	"kg/procurement/internal/vendors"
	"time"
)

type ProductID string

type Product struct {
	ID                ProductID `db:"id"`
	ProductCategoryID string    `db:"product_category_id"`
	UOMID             string    `db:"uom_id"`
	IncomeTaxID       string    `db:"income_tax_id"`
	ProductTypeID     string    `db:"product_type_id"`
	Name              string    `db:"name"`
	Description       string    `db:"description"`
	ModifiedDate      time.Time `db:"modified_date"`
	ModifiedBy        string    `db:"modified_by"`
}

type ProductVendor struct {
	Product
	vendors.Vendor
}

type GetProductsByVendorSpec struct {
	Name string `json:"name"`

	database.PaginationSpec
}
