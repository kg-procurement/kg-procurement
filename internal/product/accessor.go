package product

import (
	"context"
	"fmt"
	"kg/procurement/internal/common/database"
	"strings"
    "github.com/benbjohnson/clock"
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
    clock clock.Clock
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
		extraClauses = append(extraClauses, fmt.Sprintf("ORDER BY %s %s",
			paginationArgs.OrderBy, paginationArgs.Order))
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

func (p *postgresProductAccessor) UpdateProduct(_ context.Context, payload Product) (*Product, error) {
    now := p.clock.Now()
    
	query := `UPDATE product SET
        product_category_id = $2,
        uom_id = $3,
        income_tax_id = $4,
        product_type_id = $5,
        name = $6,
        description = $7,
        modified_date = $8,
    WHERE id = $1
    RETURNING 
        id,
        product_category_id,
        uom_id,
        income_tax_id,
        product_type_id,
        name,
        description,
        modified_date,
        modified_by
    `

    updatedProduct := &Product{}
    row := p.db.QueryRow(query,
        payload.ID,
        payload.ProductCategoryID,
        payload.UOMID,
        payload.IncomeTaxID,
        payload.ProductTypeID,
        payload.Name,
        payload.Description,
        now,
    )

    if err := row.Scan(
        &updatedProduct.ID,
        &updatedProduct.ProductCategoryID,
        &updatedProduct.UOMID,
        &updatedProduct.IncomeTaxID,
        &updatedProduct.ProductTypeID,
        &updatedProduct.Name,
        &updatedProduct.Description,
        &updatedProduct.ModifiedDate,
        &updatedProduct.ModifiedBy,
    ); err != nil {
        return nil, fmt.Errorf("failed to scan updated product: %w", err)
    }

    return updatedProduct, nil
}


func (p *postgresProductAccessor) UpdatePrice(ctx context.Context, price Price) (*Price, error) {
    now := p.clock.Now()
    query := `UPDATE price
        SET 
            purchasing_org_id = $2,
            vendor_id = $3,
            product_vendor_id = $4,
            quantity_min = $5,
            quantity_max = $6,
            quantity_uom_id = $7,
            lead_time_min = $8,
            lead_time_max = $9,
            currency_id = $10,
            price = $11,
            price_quantity = $12,
            price_uom_id = $13,
            valid_from = $14,
            valid_to = $15,
            valid_pattern_id = $16,
            area_group_id = $17,
            reference_number = $18,
            reference_date = $19,
            document_type_id = $20,
            document_id = $21,
            item_id = $22,
            term_of_payment_id = $23,
            invocation_order = $24,
            modified_date = $25,
        WHERE 
            id = $1
        RETURNING 
            id,
            purchasing_org_id,
            vendor_id,
            product_vendor_id,
            quantity_min,
            quantity_max,
            quantity_uom_id,
            lead_time_min,
            lead_time_max,
            currency_id,
            price,
            price_quantity,
            price_uom_id,
            valid_from,
            valid_to,
            valid_pattern_id,
            area_group_id,
            reference_number,
            reference_date,
            document_type_id,
            document_id,
            item_id,
            term_of_payment_id,
            invocation_order,
            modified_date,
            modified_by
    `

    updatedPrice := Price{}
    row := p.db.QueryRow(query,
        price.ID,
        price.PurchasingOrgID,
        price.VendorID,
        price.ProductVendorID,
        price.QuantityMin,
        price.QuantityMax,
        price.QuantityUOMID,
        price.LeadTimeMin,
        price.LeadTimeMax,
        price.CurrencyID,
        price.Price,
        price.PriceQuantity,
        price.PriceUOMID,
        price.ValidFrom,
        price.ValidTo,
        price.ValidPatternID,
        price.AreaGroupID,
        price.ReferenceNumber,
        price.ReferenceDate,
        price.DocumentTypeID,
        price.DocumentID,
        price.ItemID,
        price.TermOfPaymentID,
        price.InvocationOrder,
        now,
    )

    if err := row.Scan(
        &updatedPrice.ID,
        &updatedPrice.PurchasingOrgID,
        &updatedPrice.VendorID,
        &updatedPrice.ProductVendorID,
        &updatedPrice.QuantityMin,
        &updatedPrice.QuantityMax,
        &updatedPrice.QuantityUOMID,
        &updatedPrice.LeadTimeMin,
        &updatedPrice.LeadTimeMax,
        &updatedPrice.CurrencyID,
        &updatedPrice.Price,
        &updatedPrice.PriceQuantity,
        &updatedPrice.PriceUOMID,
        &updatedPrice.ValidFrom,
        &updatedPrice.ValidTo,
        &updatedPrice.ValidPatternID,
        &updatedPrice.AreaGroupID,
        &updatedPrice.ReferenceNumber,
        &updatedPrice.ReferenceDate,
        &updatedPrice.DocumentTypeID,
        &updatedPrice.DocumentID,
        &updatedPrice.ItemID,
        &updatedPrice.TermOfPaymentID,
        &updatedPrice.InvocationOrder,
        &updatedPrice.ModifiedDate,
        &updatedPrice.ModifiedBy,
    ); err != nil {
        return nil, err
    }

    return &updatedPrice, nil
}


// newPostgresProductAccessor is only accessible by the Product package
// entrypoint for other verticals should refer to the interface declared on service
func newPostgresProductAccessor(db database.DBConnector, clock clock.Clock) *postgresProductAccessor {
	return &postgresProductAccessor{
		db: db,
        clock: clock,
	}
}
