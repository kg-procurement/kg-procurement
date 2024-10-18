package account

import (
	"context"
	"kg/procurement/internal/common/database"

	"github.com/benbjohnson/clock"
)

const (
	insertAccountQuery = `
		INSERT INTO account (
			id,
			email,
			password,
			modified_date
		) VALUES (
			:id,
			:email,
			:password,
			:modified_date
		)
	`
	findAccountByEmailQuery = `
		SELECT id, email, password, modified_date, created_at
		FROM account
		WHERE email = :email
	`
)

type postgresAccountAccessor struct {
	db    database.DBConnector
	clock clock.Clock
}

func (r *postgresAccountAccessor) RegisterAccount(ctx context.Context, account Account) error {
	account.ModifiedDate = r.clock.Now()
	_, err := r.db.NamedExec(insertAccountQuery, &account)
	if err != nil {
		return err
	}
	return nil
}

func (r *postgresAccountAccessor) FindAccountByEmail(ctx context.Context, email string) (*Account, error) {
	account := &Account{}
	err := r.db.Get(account, findAccountByEmailQuery, email)
	if err != nil {
		return nil, err
	}
	return account, nil
}

// newPostgresAccountAccessor is only accessible by the Product package
// entrypoint for other verticals should refer to the interface declared on service
func newPostgresAccountAccessor(db database.DBConnector, clock clock.Clock) *postgresAccountAccessor {
	return &postgresAccountAccessor{
		db:    db,
		clock: clock,
	}
}
