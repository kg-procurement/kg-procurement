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

func (p *postgresVendorAccessor) GetAll(ctx context.Context, spec database.PaginationSpec) (*AccessorGetAllPaginationData, error) {
	args := database.BuildPaginationArgs(spec)

	dataQuery := `SELECT 
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
		ORDER BY created_at $1
		LIMIT $2
		OFFSET $3
		`

	rows, err := p.db.Query(dataQuery, args.Order, args.Limit, args.Offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

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

	countQuery := "SELECT COUNT(*) from vendor"
	totalEntries := new(int)
	row := p.db.QueryRow(countQuery)
	if err = row.Scan(&totalEntries); err != nil {
		return nil, err
	}

	metadata := database.GeneratePaginationMetadata(spec, *totalEntries)

	return &AccessorGetAllPaginationData{Vendors: vendors, Metadata: metadata}, nil
}

func (p *postgresVendorAccessor) GetByLocation(ctx context.Context, location string) ([]Vendor, error) {
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
		WHERE area_group_name = $1`

	rows, err := p.db.Query(query, location)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

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

	return vendors, nil
}

func (p *postgresVendorAccessor) GetByProductDescription(ctx context.Context, productDescription []string) ([]Vendor, error) {
	// Build the WHERE clause dynamically
	var whereClauses []string
	var args []interface{}
	for i, word := range productDescription {
		whereClauses = append(whereClauses, fmt.Sprintf("description LIKE $%d", i+1))
		args = append(args, "%"+word+"%")
	}
	whereClause := strings.Join(whereClauses, " AND ")

	// Construct the final query
	query := fmt.Sprintf(`SELECT 
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
        WHERE %s`, whereClause)

	rows, err := p.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

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

	return vendors, nil
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

func (p *postgresVendorAccessor) Put(ctx context.Context, vendor Vendor) (*Vendor, error) {
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
			modified_by = $11,
			dt = $12,
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
		vendor.ModifiedDate,
		vendor.ModifiedBy,
		vendor.Date)

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
