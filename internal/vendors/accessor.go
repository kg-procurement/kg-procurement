package vendors

import (
	"context"
	"fmt"
	"kg/procurement/internal/common/database"
	"strings"

	"github.com/benbjohnson/clock"
)

type postgresVendorAccessor struct {
	db    database.DBConnector
	clock clock.Clock
}

// GetSomeStuff is just an example
func (p *postgresVendorAccessor) GetSomeStuff(ctx context.Context) ([]string, error) {
	rows, err := p.db.Query(`SELECT name FROM users WHERE title = (?)`, "test")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		results = append(results, name)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (p *postgresVendorAccessor) GetAll(ctx context.Context, spec GetAllVendorSpec) (*AccessorGetAllPaginationData, error) {
	paginationArgs := database.BuildPaginationArgs(spec.PaginationSpec)

	// Initialize clauses and arguments
	var (
		joinClauses     []string
		whereClauses    []string
		extraClauses    []string
		args            []interface{}
		argsIndex       = 1
		extraClausesRaw = []string{
			"LIMIT $%d",
			"OFFSET $%d",
		}
	)

	// Build WHERE clause for location
	if spec.Location != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("area_group_name = $%d", argsIndex))
		args = append(args, spec.Location)
		argsIndex++
	}

	// Build JOIN and WHERE clauses for product
	if spec.Product != "" {
		joinClauses = append(joinClauses, "JOIN product_vendor pv ON pv.vendor_id = v.id")
		joinClauses = append(joinClauses, "JOIN product p ON p.id = pv.product_id")

		productNameList := strings.Fields(spec.Product)
		for _, word := range productNameList {
			whereClauses = append(whereClauses, fmt.Sprintf("p.name iLIKE $%d", argsIndex))
			args = append(args, "%"+word+"%")
			argsIndex++
		}
	}

	// Set order by default value
	if paginationArgs.OrderBy == "" {
		paginationArgs.OrderBy = "v.dt"
	} else {
		paginationArgs.OrderBy = "v." + paginationArgs.OrderBy
	}

	// Populate extra clauses
	extraClauses = append(
		extraClauses,
		fmt.Sprintf("ORDER BY %s %s", paginationArgs.OrderBy, paginationArgs.Order),
	)
	for _, clause := range extraClausesRaw {
		extraClauses = append(extraClauses, fmt.Sprintf(clause, argsIndex))
		argsIndex++
	}

	// Append pagination arguments to args
	args = append(args, paginationArgs.Limit, paginationArgs.Offset)

	// Construct the final query
	joinClause := strings.Join(joinClauses, "\n")
	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}
	extraClause := strings.Join(extraClauses, "\n")

	dataQuery := fmt.Sprintf(`
		SELECT DISTINCT
			v.id,
			v.name,
			v.description,
			v.bp_id,
			v.bp_name,
			v.rating,
			v.area_group_id,
			v.area_group_name,
			v.sap_code,
			v.modified_date,
			v.modified_by,
			v.dt
		FROM vendor v
		%s
		%s
		%s
	`, joinClause, whereClause, extraClause)

	// Execute the query
	rows, err := p.db.Query(dataQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process the results
	vendors := []Vendor{}
	for rows.Next() {
		var vendor Vendor
		err := rows.Scan(
			&vendor.ID,
			&vendor.Name,
			&vendor.Description,
			&vendor.BpID,
			&vendor.BpName,
			&vendor.Rating,
			&vendor.AreaGroupID,
			&vendor.AreaGroupName,
			&vendor.SapCode,
			&vendor.ModifiedDate,
			&vendor.ModifiedBy,
			&vendor.Date,
		)
		if err != nil {
			return nil, err
		}
		vendors = append(vendors, vendor)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Get the total count of entries
	countQuery := "SELECT COUNT(*) from vendor"
	totalEntries := new(int)
	row := p.db.QueryRow(countQuery)
	if err = row.Scan(&totalEntries); err != nil {
		return nil, err
	}

	// Generate pagination metadata
	metadata := database.GeneratePaginationMetadata(spec.PaginationSpec, *totalEntries)

	return &AccessorGetAllPaginationData{Vendors: vendors, Metadata: metadata}, nil
}

func (p *postgresVendorAccessor) GetById(ctx context.Context, id string) (*Vendor, error) {
	query := `SELECT 
		"id",
        "name",
        "description",
        "bp_id",
        "bp_name",
        "rating",
        "area_group_id",
        "area_group_name",
        "sap_code",
        "modified_date",
        "modified_by",
        "dt"
		FROM vendor 
		WHERE id = $1`

	vendor := Vendor{}
	row := p.db.QueryRow(query, id)
	if err := row.Scan(
		&vendor.ID,
		&vendor.Name,
		&vendor.Description,
		&vendor.BpID,
		&vendor.BpName,
		&vendor.Rating,
		&vendor.AreaGroupID,
		&vendor.AreaGroupName,
		&vendor.SapCode,
		&vendor.ModifiedDate,
		&vendor.ModifiedBy,
		&vendor.Date,
	); err != nil {
		return nil, err
	}

	return &vendor, nil
}

func (p *postgresVendorAccessor) UpdateDetail(ctx context.Context, vendor Vendor) (*Vendor, error) {
	now := p.clock.Now()

	// Not yet updating modified_by
	query := `UPDATE vendor
		SET 
			name = $2,
			description = $3,
			bp_id = $4,
			bp_name = $5,
			rating = $6,
			area_group_id = $7,
			area_group_name = $8,
			sap_code = $9,
			modified_date = $10,
		WHERE 
			id = $1
		RETURNING 
			id, 
			name, 
			description, 
			bp_id, 
			bp_name, 
			rating, 
			area_group_id, 
			area_group_name, 
			sap_code, 
			modified_date, 
			modified_by, 
			dt
	`

	updatedVendor := Vendor{}
	row := p.db.QueryRow(query,
		vendor.ID,
		vendor.Name,
		vendor.Description,
		vendor.BpID,
		vendor.BpName,
		vendor.Rating,
		vendor.AreaGroupID,
		vendor.AreaGroupName,
		vendor.SapCode,
		now)

	if err := row.Scan(
		&updatedVendor.ID,
		&updatedVendor.Name,
		&updatedVendor.Description,
		&updatedVendor.BpID,
		&updatedVendor.BpName,
		&updatedVendor.Rating,
		&updatedVendor.AreaGroupID,
		&updatedVendor.AreaGroupName,
		&updatedVendor.SapCode,
		&updatedVendor.ModifiedDate,
		&updatedVendor.ModifiedBy,
		&updatedVendor.Date); err != nil {
		return nil, err
	}

	return &updatedVendor, nil
}

// newPostgresVendorAccessor is only accessible by the vendor package
// entrypoint for other verticals should refer to the interface declared on service
func newPostgresVendorAccessor(db database.DBConnector, clock clock.Clock) *postgresVendorAccessor {
	return &postgresVendorAccessor{
		db:    db,
		clock: clock,
	}
}
