package product

import (
	"context"
	"fmt"
	"kg/procurement/internal/common/database"
)

const (
	getProductsByVendorQuery = `SELECT id, product_category_id, uom_id, income_tax_id, product_type_id, name, description, modified_date, modified_by FROM product_vendor WHERE vendor_id = (?)`
)

type postgresProductAccessor struct {
	db database.DBConnector
}

func (p *postgresProductAccessor) GetProductsByVendor(ctx context.Context, vendorID string) ([]Product, error) {
	rows, err := p.db.Query(getProductsByVendorQuery, vendorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := []Product{}

	for rows.Next() {
		var product Product
		err := rows.Scan(
			&product.ID,
			&product.ProductCategoryID,
			&product.UOMID,
			&product.IncomeTaxID,
			&product.ProductTypeID,
			&product.Name,
			&product.Description,
			&product.ModifiedDate,
			&product.ModifiedBy,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (p *postgresProductAccessor) UpdateProduct(_ context.Context, productID ProductID, payload Product) error {
	query := `UPDATE product SET
		name = $1 
		description = $2
		WHERE id = $3`

	_, err := p.db.Exec(query,
		payload.Name,
		payload.Description,
		productID)

	if err != nil {
		return fmt.Errorf("failed to update product with id %s: %w", productID, err)
	}

	return nil
}

// newPostgresProductAccessor is only accessible by the Product package
// entrypoint for other verticals should refer to the interface declared on service
func newPostgresProductAccessor(db database.DBConnector) *postgresProductAccessor {
	return &postgresProductAccessor{
		db: db,
	}
}
