package mailer

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/benbjohnson/clock"
	"github.com/jmoiron/sqlx"
	"github.com/onsi/gomega"
	"regexp"
	"testing"
)

func Test_newPostgresEmailStatusAccessor(t *testing.T) {
	_ = newPostgresEmailStatusAccessor(nil, nil)
}

func Test_WriteEmailStatus(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var (
			ctx     = context.Background()
			c       = setupEmailStatusAccessorTestComponent(t, WithQueryMatcher(sqlmock.QueryMatcherRegexp))
			now     = c.cmock.Now()
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
			ctx     = context.Background()
			c       = setupEmailStatusAccessorTestComponent(t)
			now     = c.cmock.Now()
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

func Test_Close(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		c := setupEmailStatusAccessorTestComponent(t)

		c.mock.ExpectClose()
		err := c.accessor.Close()
		c.g.Expect(err).ShouldNot(gomega.HaveOccurred())
	})
}

type productAccessorTestComponent struct {
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

func setupEmailStatusAccessorTestComponent(t *testing.T, opts ...Option) productAccessorTestComponent {
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

	return productAccessorTestComponent{
		g:        g,
		mock:     sqlMock,
		db:       db,
		accessor: newPostgresEmailStatusAccessor(sqlxDB, clockMock),
		cmock:    clockMock,
	}
}
