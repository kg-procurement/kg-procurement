package mailer

import (
	"context"
	"kg/procurement/internal/common/database"
	"log"

	"github.com/benbjohnson/clock"
)

const (
	insertEmailStatus = `
		INSERT INTO email_status
			(id, email_to, status, modified_date)
		VALUES
			(:id, :email_to, :status, :modified_date)
	`
)

type postgresEmailStatusAccessor struct {
	db    database.DBConnector
	clock clock.Clock
}

func (p *postgresEmailStatusAccessor) WriteEmailStatus(_ context.Context, es EmailStatus) error {
	if _, err := p.db.NamedExec(insertEmailStatus, es); err != nil {
		log.Printf("error writing email status: %v", err)
		return err
	}
	return nil
}

func (p *postgresEmailStatusAccessor) Close() error {
	return p.db.Close()
}

// newPostgresEmailStatusAccessor is only accessible by the mailer package
// entrypoint for other verticals should refer to the interface declared on service
func newPostgresEmailStatusAccessor(db database.DBConnector, clock clock.Clock) *postgresEmailStatusAccessor {
	return &postgresEmailStatusAccessor{
		db:    db,
		clock: clock,
	}
}
