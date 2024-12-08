package account

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/benbjohnson/clock"
	"github.com/jmoiron/sqlx"
	"github.com/onsi/gomega"
)

func Test_newPostgresProductAccessor(t *testing.T) {
	_ = newPostgresAccountAccessor(nil, nil)
}

func Test_RegisterAccount(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupAccountAccessorTestComponent(t, WithQueryMatcher(sqlmock.QueryMatcherRegexp))
		)
		c.cmock.Set(time.Now())
		now := c.cmock.Now()
		defer c.db.Close()

		data := Account{
			ID:           "ID",
			Email:        "a@mail.com",
			Password:     "password",
			ModifiedDate: now,
			CreatedAt:    now,
		}

		transformedQuery, args, _ := sqlx.Named(insertAccountQuery, data)
		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		c.mock.ExpectExec(transformedQuery).
			WithArgs(
				driverArgs...,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := c.accessor.RegisterAccount(ctx, data)

		c.g.Expect(err).To(gomega.BeNil())
	})

	t.Run("error - invalid query", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupAccountAccessorTestComponent(t, WithQueryMatcher(sqlmock.QueryMatcherRegexp))
		)
		c.cmock.Set(time.Now())
		now := c.cmock.Now()
		defer c.db.Close()

		data := Account{
			ID:           "ID",
			Email:        "a@mail.com",
			Password:     "password",
			ModifiedDate: now,
		}

		transformedQuery, args, _ := sqlx.Named(insertAccountQuery, data)
		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		c.mock.ExpectExec(transformedQuery).
			WithArgs(
				driverArgs...,
			).WillReturnError(sql.ErrConnDone)

		err := c.accessor.RegisterAccount(ctx, data)

		c.g.Expect(err).ToNot(gomega.BeNil())
		c.g.Expect(err).To(gomega.Equal(sql.ErrConnDone))
	})

	t.Run("error - missing fields", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupAccountAccessorTestComponent(t, WithQueryMatcher(sqlmock.QueryMatcherRegexp))
		)
		c.cmock.Set(time.Now())
		now := c.cmock.Now()
		defer c.db.Close()

		data := Account{
			ID:           "ID",
			Email:        "a@mail.com",
			ModifiedDate: now,
		}

		transformedQuery, args, _ := sqlx.Named(insertAccountQuery, data)
		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		c.mock.ExpectExec(transformedQuery).
			WithArgs(
				driverArgs...,
			).WillReturnError(sql.ErrNoRows)

		err := c.accessor.RegisterAccount(ctx, data)

		c.g.Expect(err).ToNot(gomega.BeNil())
		c.g.Expect(err).To(gomega.Equal(sql.ErrNoRows))
	})
}

func Test_FindAccountByEmail(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupAccountAccessorTestComponent(t, WithQueryMatcher(sqlmock.QueryMatcherRegexp))
		)
		defer c.db.Close()

		expectedAccount := Account{
			ID:           "ID",
			Email:        "a@mail.com",
			Password:     "hashed_password",
			ModifiedDate: time.Now(),
			CreatedAt:    time.Now(),
		}

		c.mock.ExpectQuery(findAccountByEmailQuery).
			WithArgs(expectedAccount.Email).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "modified_date", "created_at"}).
				AddRow(expectedAccount.ID, expectedAccount.Email, expectedAccount.Password, expectedAccount.ModifiedDate, expectedAccount.CreatedAt))

		account, err := c.accessor.FindAccountByEmail(ctx, expectedAccount.Email)

		c.g.Expect(err).To(gomega.BeNil())
		c.g.Expect(account).ToNot(gomega.BeNil())
		c.g.Expect(account.Email).To(gomega.Equal(expectedAccount.Email))
	})

	t.Run("error - account not found", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupAccountAccessorTestComponent(t, WithQueryMatcher(sqlmock.QueryMatcherRegexp))
		)
		defer c.db.Close()

		c.mock.ExpectQuery(findAccountByEmailQuery).
			WithArgs("nonexistent@mail.com").
			WillReturnError(sql.ErrNoRows)

		account, err := c.accessor.FindAccountByEmail(ctx, "nonexistent@mail.com")

		c.g.Expect(err).ToNot(gomega.BeNil())
		c.g.Expect(err).To(gomega.Equal(sql.ErrNoRows))
		c.g.Expect(account).To(gomega.BeNil())
	})
}

func Test_FindAccountByID(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupAccountAccessorTestComponent(t, WithQueryMatcher(sqlmock.QueryMatcherRegexp))
		)
		defer c.db.Close()

		expectedAccount := Account{
			ID:           "ID",
			Email:        "a@mail.com",
			Password:     "hashed_password",
			ModifiedDate: time.Now(),
			CreatedAt:    time.Now(),
		}

		c.mock.ExpectQuery(findAccountByIDQuery).
			WithArgs(expectedAccount.ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "modified_date", "created_at"}).
				AddRow(expectedAccount.ID, expectedAccount.Email, expectedAccount.Password, expectedAccount.ModifiedDate, expectedAccount.CreatedAt))

		account, err := c.accessor.FindAccountByID(ctx, expectedAccount.ID)

		c.g.Expect(err).To(gomega.BeNil())
		c.g.Expect(account).ToNot(gomega.BeNil())
		c.g.Expect(account.ID).To(gomega.Equal(expectedAccount.ID))
		c.g.Expect(account.Email).To(gomega.Equal(expectedAccount.Email))
	})

	t.Run("error - account not found", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupAccountAccessorTestComponent(t, WithQueryMatcher(sqlmock.QueryMatcherRegexp))
		)
		defer c.db.Close()

		c.mock.ExpectQuery(findAccountByIDQuery).
			WithArgs("nonexistentID").
			WillReturnError(sql.ErrNoRows)

		account, err := c.accessor.FindAccountByID(ctx, "nonexistentID")

		c.g.Expect(err).ToNot(gomega.BeNil())
		c.g.Expect(err).To(gomega.Equal(sql.ErrNoRows))
		c.g.Expect(account).To(gomega.BeNil())
	})

	t.Run("error - query execution failure", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupAccountAccessorTestComponent(t, WithQueryMatcher(sqlmock.QueryMatcherRegexp))
		)
		defer c.db.Close()

		c.mock.ExpectQuery(findAccountByIDQuery).
			WithArgs("ID").
			WillReturnError(sql.ErrConnDone)

		account, err := c.accessor.FindAccountByID(ctx, "ID")

		c.g.Expect(err).ToNot(gomega.BeNil())
		c.g.Expect(err).To(gomega.Equal(sql.ErrConnDone))
		c.g.Expect(account).To(gomega.BeNil())
	})
}


type accountAccessorTestComponent struct {
	g        *gomega.WithT
	mock     sqlmock.Sqlmock
	db       *sql.DB
	accessor *postgresAccountAccessor
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

func setupAccountAccessorTestComponent(t *testing.T, opts ...Option) accountAccessorTestComponent {
	options := setupOptions{
		queryMatcher: sqlmock.QueryMatcherEqual,
	}

	for _, opt := range opts {
		opt(&options)
	}

	g := gomega.NewWithT(t)
	db, sqlMock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	clockMock := clock.NewMock()

	return accountAccessorTestComponent{
		g:        g,
		mock:     sqlMock,
		db:       db,
		accessor: newPostgresAccountAccessor(sqlxDB, clockMock),
		cmock:    clockMock,
	}
}
