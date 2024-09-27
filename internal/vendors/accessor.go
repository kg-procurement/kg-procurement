package vendors

import (
	"context"
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

// newPostgresVendorAccessor is only accessible by the vendor package
// entrypoint for other verticals should refer to the interface declared on service
func newPostgresVendorAccessor(db database.DBConnector) *postgresVendorAccessor {
	return &postgresVendorAccessor{
		db: db,
	}
}
