package mailer

import (
	"context"
	"fmt"
	"kg/procurement/cmd/utils"
	"kg/procurement/internal/common/database"
	"log"
	"strings"

	"github.com/benbjohnson/clock"
)

const (
	insertEmailStatus = `
		INSERT INTO email_status
			(id, email_to, status, modified_date)
		VALUES
			(:id, :email_to, :status, :modified_date)
	`
	updateEmailStatus = `
		UPDATE email_status
		SET status = :status, modified_date = :modified_date
		WHERE id = :id
		RETURNING id, email_to, status, modified_date
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

func (p *postgresEmailStatusAccessor) UpdateEmailStatus(_ context.Context, es EmailStatus) (*EmailStatus, error) {
	es.ModifiedDate = p.clock.Now()
	var updatedEmailStatus EmailStatus
	rows, err := p.db.NamedQuery(updateEmailStatus, es)
	if err != nil {
		log.Printf("error updating email status: %v", err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(&updatedEmailStatus); err != nil {
			log.Printf("error scanning updated email status: %v", err)
			return nil, err
		}
	}
	return &updatedEmailStatus, nil
}

func (p *postgresEmailStatusAccessor) GetAll(ctx context.Context, spec GetAllEmailStatusSpec) (*AccessorGetAllPaginationData, error) {
	paginationArgs := database.BuildPaginationArgs(spec.PaginationSpec)

	var (
		whereClauses    []string
		joinClauses     []string
		extraClauses    []string
		args            []interface{}
		countArgs       []interface{}
		argsIndex       = 1
		extraClausesRaw = []string{
			"LIMIT $%d",
			"OFFSET $%d",
		}
	)

	// Build WHERE clauses
	if spec.EmailTo != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("es.email_to ILIKE $%d", argsIndex))
		argValue := "%" + spec.EmailTo + "%"
		args = append(args, argValue)
		countArgs = append(countArgs, argValue)
		argsIndex++
	}

	// Set order by default value
	if paginationArgs.OrderBy == "" {
		paginationArgs.OrderBy = "es.modified_date"
	} else {
		paginationArgs.OrderBy = "es." + paginationArgs.OrderBy
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

	// Construct the WHERE clause
	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Construct the final query
	joinClause := strings.Join(joinClauses, "\n")
	extraClause := strings.Join(extraClauses, "\n")

	dataQuery := fmt.Sprintf(`
        SELECT DISTINCT
            es.id,
            es.email_to,
            es.status,
            es.modified_date
        FROM email_status es
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
	emailStatus := []EmailStatus{}
	for rows.Next() {
		var es EmailStatus
		if err := rows.StructScan(&es); err != nil {
			return nil, err
		}
		emailStatus = append(emailStatus, es)
	}
	if err := rows.Err(); err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	// Get the total count of entries
	countQuery := "SELECT COUNT(*) from email_status es"
	if whereClause != "" {
		countQuery += " " + whereClause
	}
	totalEntries := 0
	err = p.db.QueryRow(countQuery, countArgs...).Scan(&totalEntries)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	// Generate pagination metadata
	metadata := database.GeneratePaginationMetadata(spec.PaginationSpec, totalEntries)

	return &AccessorGetAllPaginationData{EmailStatus: emailStatus, Metadata: metadata}, nil
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
