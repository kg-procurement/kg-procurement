package vendors

import (
	"context"
	"fmt"
	"kg/procurement/internal/common/database"
)

type postgresVendorAccessor struct {
	db database.DBConnector
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

func (p *postgresVendorAccessor) GetAll(ctx context.Context) ([]Vendor, error) {
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
		FROM vendor`

	rows, err := p.db.Query(query)
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
			&vendor.BpId,
			&vendor.BpName,
			&vendor.Rating,
			&vendor.AreaGroupId,
			&vendor.AreaGroupName,
			&vendor.SapCode,
			&vendor.ModifiedDate,
			&vendor.ModifiedBy,
			&vendor.Date,
		)
		if err != nil {
			return nil, fmt.Errorf("Failed while scanning row: %w", err)
		}
		vendors = append(vendors, vendor)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return vendors, nil
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
			&vendor.BpId,
			&vendor.BpName,
			&vendor.Rating,
			&vendor.AreaGroupId,
			&vendor.AreaGroupName,
			&vendor.SapCode,
			&vendor.ModifiedDate,
			&vendor.ModifiedBy,
			&vendor.Date,
		)
		if err != nil {
			return nil, fmt.Errorf("Failed while scanning row: %w", err)
		}
		vendors = append(vendors, vendor)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return vendors, nil
}

// newPostgresVendorAccessor is only accessible by the vendor package
// entrypoint for other verticals should refer to the interface declared on service
func newPostgresVendorAccessor(db database.DBConnector) *postgresVendorAccessor {
	return &postgresVendorAccessor{
		db: db,
	}
}
