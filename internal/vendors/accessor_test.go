package vendors

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/onsi/gomega"
)

func Test_postgresVendorAccessor_GetSomeStuff(t *testing.T) {
	t.Parallel()

	var (
		accessor *postgresVendorAccessor
		mock     sqlmock.Sqlmock
	)

	setup := func(t *testing.T) (*gomega.GomegaWithT, *sql.DB) {
		g := gomega.NewWithT(t)
		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			log.Fatal("error initializing mock:", err)
		}

		accessor = newPostgresVendorAccessor(db)
		mock = sqlMock

		return g, db
	}

	t.Run("success", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice")

		mock.ExpectQuery(`SELECT name FROM users WHERE title = (?)`).
			WithArgs("test").
			WillReturnRows(rows)

		ctx := context.Background()
		res, err := accessor.GetSomeStuff(ctx)

		g.Expect(err).To(gomega.BeNil())
		g.Expect(res).To(gomega.Equal([]string{"Alice"}))
	})

	t.Run("error on query", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		mock.ExpectQuery(`SELECT name FROM users WHERE title = (?)`).
			WithArgs("test").
			WillReturnError(errors.New("some error"))

		ctx := context.Background()
		res, err := accessor.GetSomeStuff(ctx)

		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(res).To(gomega.BeNil())
	})

	t.Run("error on row scan", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow(nil)

		mock.ExpectQuery(`SELECT name FROM users WHERE title = (?)`).
			WithArgs("test").
			WillReturnRows(rows)

		ctx := context.Background()
		res, err := accessor.GetSomeStuff(ctx)

		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(res).To(gomega.BeNil())
	})
}

func TestVendorAccessor_GetAll(t *testing.T) {
	t.Parallel()

	var (
		accessor *postgresVendorAccessor
		mock     sqlmock.Sqlmock
	)

	setup := func(t *testing.T) (*gomega.GomegaWithT, *sql.DB) {
		g := gomega.NewWithT(t)
		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			log.Fatal("error initializing mock:", err)
		}

		accessor = newPostgresVendorAccessor(db)
		mock = sqlMock

		return g, db
	}

	t.Run("success", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()
		sampleData := []string{
			"id",
			"name",
			"bp_id",
			"bp_name",
			"rating",
			"area_group_id",
			"area_group_name",
			"sap_code",
			"modified_date",
			"modified_by",
			"dt",
		}

		rows := sqlmock.NewRows(sampleData).
			AddRow(
				1,
				"name",
				1,
				"bp_name",
				1,
				1,
				"group_name",
				"sap_code",
				time.Now(),
				1,
				time.Now(),
			)

		mock.ExpectQuery(`SELECT * FROM vendor`).
			WillReturnRows(rows)

		ctx := context.Background()
		res, err := accessor.GetAll(ctx)
		g.Expect(err).To(gomega.BeNil())
		g.Expect(res).To(gomega.Equal(sampleData))
	})
}
