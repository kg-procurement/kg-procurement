package inventory

import (
	"context"
	"fmt"
	"kg/procurement/internal/common/database"
)

type postgresInventoryAccessor struct {
	db database.DBConnector
}

func (p *postgresInventoryAccessor) UpdateInventory(ctx context.Context, id string, updatedItem Item) error {
	query := `UPDATE inventory SET 
		name = $1, 
		description = $2, 
		price = $3
		WHERE id = $4`

	_, err := p.db.Exec(query, 
		updatedItem.Name, 
		updatedItem.Description, 
		updatedItem.Price, 
		id)
	
	if err != nil {
		return fmt.Errorf("failed to update inventory item with id %s: %w", id, err)
	}
	
	return nil
}

// newPostgresInventoryAccessor is only accessible by the Inventory package
// entrypoint for other verticals should refer to the interface declared on service
func newPostgresInventoryAccessor(db database.DBConnector) *postgresInventoryAccessor {
	return &postgresInventoryAccessor{
		db: db,
	}
}
