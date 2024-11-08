package mailer

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/benbjohnson/clock"
	"github.com/jmoiron/sqlx"
	"github.com/onsi/gomega"
	"kg/procurement/internal/common/database"
	"log"
	"regexp"
	"testing"
	"time"
)

func Test_newPostgresEmailStatusAccessor(t *testing.T) {
	_ = newPostgresEmailStatusAccessor(nil, nil)
}

func Test_WriteEmailStatus(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var (
			ctx         = context.Background()
			c           = setupEmailStatusAccessorTestComponent(t, WithQueryMatcher(sqlmock.QueryMatcherRegexp))
			now         = c.cmock.Now()
			emailStatus = EmailStatus{
				ID:           "123",
				EmailTo:      "email@email.com",
				Status:       "sent",
				ModifiedDate: now,
			}
		)

		transformedQuery, args, _ := sqlx.Named(insertEmailStatus, emailStatus)
		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		c.mock.ExpectExec(regexp.QuoteMeta(transformedQuery)).WithArgs(
			driverArgs...,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err := c.accessor.WriteEmailStatus(ctx, emailStatus)
		c.g.Expect(err).Should(gomega.BeNil())
	})

	t.Run("returns error on db failure", func(t *testing.T) {
		var (
			ctx         = context.Background()
			c           = setupEmailStatusAccessorTestComponent(t)
			now         = c.cmock.Now()
			emailStatus = EmailStatus{
				ID:           "123",
				EmailTo:      "email@email.com",
				Status:       "sent",
				ModifiedDate: now,
			}
		)

		transformedQuery, args, _ := sqlx.Named(insertEmailStatus, emailStatus)
		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		c.mock.ExpectExec(regexp.QuoteMeta(transformedQuery)).WithArgs(
			driverArgs...,
		).WillReturnError(sql.ErrConnDone)

		err := c.accessor.WriteEmailStatus(ctx, emailStatus)
		c.g.Expect(err).Should(gomega.HaveOccurred())
	})
}

func Test_GetAll(t *testing.T) {
	t.Parallel()

	var (
		accessor *postgresEmailStatusAccessor
		mock     sqlmock.Sqlmock
	)

	setup := func(t *testing.T) (*gomega.WithT, *sql.DB) {
		g := gomega.NewWithT(t)
		realClock := clock.New()
		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			log.Fatal("error initializing mock:", err)
		}
		sqlxDB := sqlx.NewDb(db, "sqlmock")

		accessor = newPostgresEmailStatusAccessor(sqlxDB, realClock)
		mock = sqlMock

		return g, db
	}

	emailStatusFields := []string{
		"id",
		"email_to",
		"status",
		"modified_date",
	}

	dataQuery := `
		SELECT DISTINCT
			es.id,
			es.email_to,
			es.status,
			es.modified_date
		FROM email_status es
		ORDER BY es.modified_date DESC
		LIMIT $1
		OFFSET $2
	`

	countQuery := "SELECT COUNT(*) from email_status"

	fixedTime := time.Date(2024, time.September, 23, 12, 30, 0, 0, time.UTC)

	spec := GetAllEmailStatusSpec{
		PaginationSpec: database.PaginationSpec{
			Order: "DESC",
			Limit: 10,
			Page:  1,
		},
	}

	t.Run("success", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(emailStatusFields).
			AddRow(
				"1",
				"test@example.com",
				"sent",
				fixedTime,
			)

		args := database.BuildPaginationArgs(spec.PaginationSpec)

		mock.ExpectQuery(dataQuery).
			WithArgs(args.Limit, args.Offset).
			WillReturnRows(rows)

		totalRows := sqlmock.NewRows([]string{"count"}).AddRow(1)

		mock.ExpectQuery(countQuery).WillReturnRows(totalRows)

		ctx := context.Background()
		res, err := accessor.GetAll(ctx, spec)

		emailStatusExpectation := []EmailStatus{{
			ID:           "1",
			EmailTo:      "test@example.com",
			Status:       "sent",
			ModifiedDate: fixedTime,
		}}

		expectation := &AccessorGetAllPaginationData{
			EmailStatus: emailStatusExpectation,
			Metadata:    res.Metadata,
		}

		g.Expect(err).To(gomega.BeNil())
		g.Expect(res).To(gomega.Equal(expectation))
	})
}

func Test_Close(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		c := setupEmailStatusAccessorTestComponent(t)

		c.mock.ExpectClose()
		err := c.accessor.Close()
		c.g.Expect(err).ShouldNot(gomega.HaveOccurred())
	})
}

type emailStatusAccessorTestComponent struct {
	g        *gomega.WithT
	mock     sqlmock.Sqlmock
	db       *sql.DB
	accessor *postgresEmailStatusAccessor
	cmock    *clock.Mock
}

type setupOptions struct {
	queryMatcher sqlmock.QueryMatcher
}

func WithQueryMatcher(matcher sqlmock.QueryMatcher) Option {
	return func(o *setupOptions) {
		o.queryMatcher = matcher
	}
}

type Option func(*setupOptions)

func setupEmailStatusAccessorTestComponent(t *testing.T, opts ...Option) emailStatusAccessorTestComponent {
	options := setupOptions{
		queryMatcher: sqlmock.QueryMatcherEqual,
	}

	for _, opt := range opts {
		opt(&options)
	}

	g := gomega.NewWithT(t)
	db, sqlMock, _ := sqlmock.New(sqlmock.QueryMatcherOption(options.queryMatcher))
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	clockMock := clock.NewMock()

	return emailStatusAccessorTestComponent{
		g:        g,
		mock:     sqlMock,
		db:       db,
		accessor: newPostgresEmailStatusAccessor(sqlxDB, clockMock),
		cmock:    clockMock,
	}
}
