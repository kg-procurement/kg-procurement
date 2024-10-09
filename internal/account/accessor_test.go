package account

import (
	"context"
	"database/sql"
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

	accountFields := []string{
		"id",
		"email",
		"password",
		"created_at",
	}

	t.Run("success", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupProductAccessorTestComponent(t)
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

		c.mock.ExpectQuery(insertAccountQuery).
			WithArgs(
				"ID",
				"a@mail.com",
				"password",
				now,
			).WillReturnRows(sqlmock.NewRows(accountFields))

		err := c.accessor.RegisterAccount(ctx, data)

		c.g.Expect(err).To(gomega.BeNil())
	})

	t.Run("error - invalid query", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupProductAccessorTestComponent(t)
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

		c.mock.ExpectQuery(insertAccountQuery).
			WithArgs(
				"ID",
				"a@mail.com",
				"password",
				now,
			).WillReturnError(sql.ErrConnDone)

		err := c.accessor.RegisterAccount(ctx, data)

		c.g.Expect(err).ToNot(gomega.BeNil())
		c.g.Expect(err).To(gomega.Equal(sql.ErrConnDone))
	})

	t.Run("error - missing fields", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupProductAccessorTestComponent(t)
		)
		c.cmock.Set(time.Now())
		now := c.cmock.Now()
		defer c.db.Close()

		data := Account{
			ID:           "ID",
			Email:        "a@mail.com",
			ModifiedDate: now,
		}

		c.mock.ExpectQuery(insertAccountQuery).
			WithArgs(
				"ID",
				"a@mail.com",
				sqlmock.AnyArg(),
				now,
			).WillReturnError(sql.ErrNoRows)

		err := c.accessor.RegisterAccount(ctx, data)

		c.g.Expect(err).ToNot(gomega.BeNil())
		c.g.Expect(err).To(gomega.Equal(sql.ErrNoRows))
	})
}

type productAccessorTestComponent struct {
	g        *gomega.WithT
	mock     sqlmock.Sqlmock
	db       *sql.DB
	accessor *postgresAccountAccessor
	cmock    *clock.Mock
}

func setupProductAccessorTestComponent(t *testing.T) productAccessorTestComponent {
	g := gomega.NewWithT(t)
	db, sqlMock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	clockMock := clock.NewMock()

	return productAccessorTestComponent{
		g:        g,
		mock:     sqlMock,
		db:       db,
		accessor: newPostgresAccountAccessor(sqlxDB, clockMock),
		cmock:    clockMock,
	}
}
