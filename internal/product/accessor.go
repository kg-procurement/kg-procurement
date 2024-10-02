package product

import (
	"context"
	"errors"
	"fmt"
	"kg/procurement/internal/common/database"
)

const (
	getProductsByVendorQuery = `SELECT id, product_category_id, uom_id, income_tax_id, product_type_id, name, description, modified_date, modified_by FROM product_vendor WHERE vendor_id = (?)`
)

type postgresProductAccessor struct {
	db database.DBConnector
}

func (p *postgresProductAccessor) GetProductsByVendor(_ context.Context, _ string) ([]Product, error) {
	return nil, errors.New("implement me")
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
