package vendors

import (
	"context"
	"kg/procurement/internal/common/database"
)

type postgresVendorAccessor struct {
	db *database.Conn
}

func (p *postgresVendorAccessor) GetSomeStuff(_ context.Context) error {
	// db queries here
	return nil
}

// newPostgresVendorAccessor is only accessible by the vendor package
// entrypoint for other verticals should refer to the interface declared on service
func newPostgresVendorAccessor(db *database.Conn) *postgresVendorAccessor {
	return &postgresVendorAccessor{
		db: db,
	}
}
