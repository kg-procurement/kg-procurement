package account

import (
	"context"
	"kg/procurement/internal/common/database"

	"github.com/benbjohnson/clock"
)

const (
	insertAccountQuery = `
		INSERT INTO accounts (
			id,
			email,
			password,
			modified_date
		) VALUES (
			$1,
			$2,
			$3,
			$4
		)
	`
)

type postgresAccountAccessor struct {
	db    database.DBConnector
	clock clock.Clock
}

func (r *postgresAccountAccessor) RegisterAccount(ctx context.Context, account Account) error {
	return nil
}

// newPostgresAccountAccessor is only accessible by the Product package
// entrypoint for other verticals should refer to the interface declared on service
func newPostgresAccountAccessor(db database.DBConnector, clock clock.Clock) *postgresAccountAccessor {
	return &postgresAccountAccessor{
		db:    db,
		clock: clock,
	}
}
