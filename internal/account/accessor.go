package account

import (
	"context"
	"kg/procurement/cmd/utils"
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
		WHERE email = $1
	`

	findAccountByIDQuery = `
		SELECT id, email, password, modified_date, created_at
		FROM account
		WHERE id = $1
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
		utils.Logger.Error(err.Error())
		return err
	}
	return nil
}

func (r *postgresAccountAccessor) FindAccountByEmail(ctx context.Context, email string) (*Account, error) {
	account := &Account{}
	err := r.db.Get(account, findAccountByEmailQuery, email)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}
	return account, nil
}

func (r *postgresAccountAccessor) FindAccountByID(ctx context.Context, id string) (*Account, error) {
	account := &Account{}
	err := r.db.Get(account, findAccountByIDQuery, id)
	if err != nil {
		utils.Logger.Error(err.Error())
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
