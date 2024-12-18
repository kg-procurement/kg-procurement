package vendors

import (
	"context"
	"fmt"
	"kg/procurement/cmd/utils"
	"kg/procurement/internal/common/database"
	"strings"

	"github.com/benbjohnson/clock"
	"github.com/jmoiron/sqlx"
)

type postgresVendorAccessor struct {
	db    database.DBConnector
	clock clock.Clock
}

const (
	insertVendor = `
		INSERT INTO vendor
			(id, name, email, description, bp_id, bp_name, rating, area_group_id, area_group_name, sap_code, modified_date, modified_by, dt)
		VALUES 
			(:id, :name, :email, :description, :bp_id, :bp_name, :rating, :area_group_id, :area_group_name, :sap_code, :modified_date, :modified_by, :dt)
	`
	getBulkByID          = `SELECT * FROM vendor WHERE id IN (?)`
	getAllLocationsQuery = `SELECT DISTINCT area_group_name FROM vendor`
	getBulkByProductName = `
		SELECT
			v.*
		FROM
			vendor v
		JOIN 
			price pp ON pp.vendor_id = v.id
		JOIN
			product_vendor pv ON pv.id = pp.product_vendor_id
		JOIN 
			product p ON p.id = pv.product_id
		WHERE
			p.name = :product_name
	`
	createEvaluationQuery = `
		INSERT INTO vendor_evaluation
			(id, vendor_id, kesesuaian_produk, kualitas_produk, ketepatan_waktu_pengiriman, kompetitifitas_harga, responsivitas_kemampuan_komunikasi, kemampuan_dalam_menangani_masalah, kelengkapan_barang, harga, term_of_payment, reputasi, ketersediaan_barang, kualitas_layanan_after_services, modified_date)
		VALUES
			(:id, :vendor_id, :kesesuaian_produk, :kualitas_produk, :ketepatan_waktu_pengiriman, :kompetitifitas_harga, :responsivitas_kemampuan_komunikasi, :kemampuan_dalam_menangani_masalah, :kelengkapan_barang, :harga, :term_of_payment, :reputasi, :ketersediaan_barang, :kualitas_layanan_after_services, :modified_date)
	`
)

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
		whereClauses = append(whereClauses, fmt.Sprintf("v.area_group_name = $%d", argsIndex))
		args = append(args, spec.Location)
		argsIndex++
	}

	// Build JOIN and WHERE clauses for product
	if spec.Product != "" {
		joinClauses = append(joinClauses, "JOIN price pr ON pr.vendor_id = v.id")
		joinClauses = append(joinClauses, "JOIN product_vendor pv ON pv.id = pr.product_vendor_id")
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
	rows, err := p.db.Queryx(dataQuery, args...)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	// Process the results
	vendors := []Vendor{}
	for rows.Next() {
		var vendor Vendor
		if err := rows.StructScan(&vendor); err != nil {
			return nil, err
		}
		vendors = append(vendors, vendor)
	}
	if err := rows.Err(); err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	// Get the total count of entries
	countQuery := "SELECT COUNT(*) from vendor"
	totalEntries := new(int)
	row := p.db.QueryRow(countQuery)
	if err = row.Scan(&totalEntries); err != nil {
		utils.Logger.Error(err.Error())
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
	row := p.db.QueryRowx(query, id)
	if err := row.StructScan(&vendor); err != nil {
		utils.Logger.Error(err.Error())
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
			modified_date = $10
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

	updatedVendor := &Vendor{}
	row := p.db.QueryRowx(query,
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

	if err := row.StructScan(updatedVendor); err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	return updatedVendor, nil
}

func (p *postgresVendorAccessor) GetAllLocations(ctx context.Context) ([]string, error) {
	rows, err := p.db.Query(getAllLocationsQuery)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			utils.Logger.Error(err.Error())
			return nil, err
		}
		results = append(results, name)
	}

	if err := rows.Err(); err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	return results, nil
}

func (p *postgresVendorAccessor) BulkGetByIDs(_ context.Context, ids []string) ([]Vendor, error) {
	query, args, err := sqlx.In(getBulkByID, ids)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	query = p.db.Rebind(query)
	rows, err := p.db.Queryx(query, args...)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	res := []Vendor{}
	for rows.Next() {
		var vendor Vendor
		if err := rows.StructScan(&vendor); err != nil {
			utils.Logger.Error(err.Error())
			return nil, err
		}
		res = append(res, vendor)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (p *postgresVendorAccessor) writeVendor(ctx context.Context, vendor Vendor) error {
	if _, err := p.db.NamedExec(insertVendor, vendor); err != nil {
		utils.Logger.Error(err.Error())
		return err
	}
	return nil
}

func (p *postgresVendorAccessor) BulkGetByProductName(_ context.Context, productName string) ([]Vendor, error) {
	query := getBulkByProductName

	rows, err := p.db.NamedQuery(query, map[string]interface{}{
		"product_name": productName,
	})

	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	vendors := []Vendor{}
	for rows.Next() {
		var vendor Vendor
		if err = rows.StructScan(&vendor); err != nil {
			utils.Logger.Error(err.Error())
			return nil, err
		}
		vendors = append(vendors, vendor)
	}
	if err := rows.Err(); err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	return vendors, nil
}

func (p *postgresVendorAccessor) CreateEvaluation(ctx context.Context, evaluation *VendorEvaluation) (*VendorEvaluation, error) {
	if _, err := p.db.NamedExec(createEvaluationQuery, evaluation); err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	return evaluation, nil
}

func (p *postgresVendorAccessor) Close() error {
	return p.db.Close()
}

// newPostgresVendorAccessor is only accessible by the vendor package
// entrypoint for other verticals should refer to the interface declared on service
func newPostgresVendorAccessor(db database.DBConnector, clock clock.Clock) *postgresVendorAccessor {
	return &postgresVendorAccessor{
		db:    db,
		clock: clock,
	}
}
