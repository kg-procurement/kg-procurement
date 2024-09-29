package inventory

import (
	"context"
	"kg/procurement/internal/common/database"
)

type postgresInventoryAccessor struct {
	db database.DBConnector
}

func (p *postgresInventoryAccessor) UpdateInventory(ctx context.Context, id string, updatedItem Item) error {
	return nil
}

// newPostgresInventoryAccessor is only accessible by the Inventory package
// entrypoint for other verticals should refer to the interface declared on service
func newPostgresInventoryAccessor(db database.DBConnector) *postgresInventoryAccessor {
	return &postgresInventoryAccessor{
		db: db,
	}
}
