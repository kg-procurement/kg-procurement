package product

import (
	"context"
	"fmt"
	"kg/procurement/internal/common/database"
	"strings"
)

const (
	getProductsByVendorQuery = `
		SELECT 
			p.id, 
			p.product_category_id, 
			p.uom_id, 
			p.income_tax_id, 
			p.product_type_id, 
			p.name, 
			p.description, 
			p.modified_date, 
			p.modified_by 
		FROM 
			product p
		JOIN 
			product_vendor pv ON pv.product_id = p.id
		WHERE 
			pv.vendor_id = $1
	`
)

type postgresProductAccessor struct {
	db database.DBConnector
}

func (p *postgresProductAccessor) GetProductsByVendor(
	_ context.Context, vendorID string, spec GetProductsByVendorSpec) ([]Product, error) {
	paginationArgs := database.BuildPaginationArgs(spec.PaginationSpec)

	// Initialize clauses and arguments
	var (
		whereClauses []string
		extraClauses []string
		args         = []interface{}{vendorID}
		argsIndex    = 2 // start at 2 because the query already have $1
	)

	// Build WHERE clauses for product
	if spec.Name != "" {
		productNameList := strings.Fields(spec.Name)
		for _, word := range productNameList {
			whereClauses = append(whereClauses, fmt.Sprintf("p.name iLIKE $%d", argsIndex))
			args = append(args, "%"+word+"%")
			argsIndex++
		}
	}

	// Build extra clauses
	if paginationArgs.OrderBy != "" {
		extraClauses = append(extraClauses, fmt.Sprintf("ORDER BY %s", paginationArgs.OrderBy))
	}

	// Build the query
	query := getProductsByVendorQuery
	if len(whereClauses) > 0 {
		query += " AND " + strings.Join(whereClauses, " AND ")
	}
	if len(extraClauses) > 0 {
		query += " " + strings.Join(extraClauses, " ")
	}

	rows, err := p.db.Query(query, args...)
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
