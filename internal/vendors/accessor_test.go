package vendors

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/onsi/gomega"
)

func Test_newPostgresVendorAccessor(t *testing.T) {
	_ = newPostgresVendorAccessor(nil)
}

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

	sampleData := []string{
		"id",
		"name",
		"description",
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

	query := `SELECT 
			"id",
			"name",
			"description",
			"bp_id",
			"bp_name",
			"rating",
			"area_group_id",
			"area_group_name",
			"sap_code",
			"modified_date",
			"modified_by",
			"dt" 
			FROM vendor`

	fixedTime := time.Date(2024, time.September, 27, 12, 30, 0, 0, time.UTC)

	t.Run("success", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(sampleData).
			AddRow(
				"1",
				"name",
				"description",
				1,
				"bp_name",
				1,
				1,
				"group_name",
				"sap_code",
				fixedTime,
				1,
				fixedTime,
			)

		mock.ExpectQuery(query).
			WillReturnRows(rows)

		ctx := context.Background()
		res, err := accessor.GetAll(ctx)

		expectation := []Vendor{{
			ID:            "1",
			Name:          "name",
			Description:   "description",
			BpID:          "1",
			BpName:        "bp_name",
			Rating:        1,
			AreaGroupID:   "1",
			AreaGroupName: "group_name",
			SapCode:       "sap_code",
			ModifiedDate:  fixedTime,
			ModifiedBy:    1,
			Date:          fixedTime,
		}}

		g.Expect(err).To(gomega.BeNil())
		g.Expect(res).To(gomega.Equal(expectation))
	})

	t.Run("success on empty result", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(sampleData)

		mock.ExpectQuery(query).
			WillReturnRows(rows)

		ctx := context.Background()
		res, err := accessor.GetAll(ctx)
		g.Expect(err).To(gomega.BeNil())
		g.Expect(res).To(gomega.Equal([]Vendor{}))
	})

	t.Run("error on scanning row", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(sampleData).AddRow(
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)

		mock.ExpectQuery(query).
			WillReturnRows(rows)

		ctx := context.Background()
		res, err := accessor.GetAll(ctx)

		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(res).To(gomega.BeNil())
	})

	t.Run("error on executing query", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(sampleData).
			AddRow(
				"1",
				"name",
				"description",
				1,
				"bp_name",
				1,
				1,
				"group_name",
				"sap_code",
				fixedTime,
				1,
				fixedTime,
			)

		wrongQuery := `SELECT 
			"name",
			"description",
			"bp_id",
			"bp_name",
			"rating",
			"area_group_id",
			"area_group_name",
			"sap_code",
			"modified_date",
			"modified_by",
			"dt" 
			FROM user`

		mock.ExpectQuery(wrongQuery).
			WillReturnRows(rows)

		ctx := context.Background()
		res, err := accessor.GetAll(ctx)

		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(res).To(gomega.BeNil())
	})

	t.Run("error while iterating rows", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(sampleData).
			AddRow(
				"1",
				"name",
				"description",
				1,
				"bp_name",
				1,
				1,
				"group_name",
				"sap_code",
				fixedTime,
				1,
				fixedTime,
			).AddRow(
			"1",
			"name",
			"description",
			1,
			"bp_name",
			1,
			1,
			"group_name",
			"sap_code",
			fixedTime,
			1,
			fixedTime,
		).RowError(1, fmt.Errorf("row error"))

		mock.ExpectQuery(query).
			WillReturnRows(rows)

		ctx := context.Background()
		res, err := accessor.GetAll(ctx)

		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(res).To(gomega.BeNil())
	})
}

func TestVendorAccessor_GetByLocation(t *testing.T) {
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

	sampleData := []string{
		"id",
		"name",
		"description",
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

	query := `SELECT 
	"id",
	"name",
	"description",
	"bp_id",
	"bp_name",
	"rating",
	"area_group_id",
	"area_group_name",
	"sap_code",
	"modified_date",
	"modified_by",
	"dt" 
	FROM vendor
	WHERE area_group_name = $1`

	fixedTime := time.Date(2024, time.September, 27, 12, 30, 0, 0, time.UTC)

	location := "Indonesia"

	t.Run("success", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(sampleData).
			AddRow(
				"1",
				"name",
				"description",
				"1",
				"bp_name",
				1,
				"1",
				location,
				"sap_code",
				fixedTime,
				1,
				fixedTime,
			)

		mock.ExpectQuery(query).
			WithArgs(location).
			WillReturnRows(rows)

		ctx := context.Background()
		res, err := accessor.GetByLocation(ctx, location)

		expectation := []Vendor{{
			ID:            "1",
			Name:          "name",
			Description:   "description",
			BpID:          "1",
			BpName:        "bp_name",
			Rating:        1,
			AreaGroupID:   "1",
			AreaGroupName: location,
			SapCode:       "sap_code",
			ModifiedDate:  fixedTime,
			ModifiedBy:    1,
			Date:          fixedTime,
		}}

		g.Expect(err).To(gomega.BeNil())
		g.Expect(res).To(gomega.Equal(expectation))
	})

	t.Run("success on empty result", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(sampleData)

		mock.ExpectQuery(query).
			WithArgs(location).
			WillReturnRows(rows)

		ctx := context.Background()
		res, err := accessor.GetByLocation(ctx, location)
		g.Expect(err).To(gomega.BeNil())
		g.Expect(res).To(gomega.Equal([]Vendor{}))
	})

	t.Run("error on scanning row", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(sampleData).AddRow(
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)

		mock.ExpectQuery(query).
			WithArgs(location).
			WillReturnRows(rows)

		ctx := context.Background()
		res, err := accessor.GetByLocation(ctx, location)

		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(res).To(gomega.BeNil())
	})

	t.Run("error on executing query", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(sampleData).
			AddRow(
				"1",
				"name",
				"description",
				"1",
				"bp_name",
				1,
				"1",
				location,
				"sap_code",
				fixedTime,
				1,
				fixedTime,
			)

		wrongQuery := `SELECT 
			"id",
			"name",
			"description",
			"bp_id",
			"bp_name",
			"rating",
			"area_group_id",
			"area_group_name",
			"sap_code",
			"modified_date",
			"modified_by",
			"dt" 
			FROM vendor
			WHERE description = $1`

		mock.ExpectQuery(wrongQuery).
			WillReturnRows(rows)

		ctx := context.Background()
		res, err := accessor.GetByLocation(ctx, location)

		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(res).To(gomega.BeNil())
	})

	t.Run("error while iterating rows", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(sampleData).
			AddRow(
				"1",
				"name",
				"description",
				"1",
				"bp_name",
				1,
				"1",
				location,
				"sap_code",
				fixedTime,
				1,
				fixedTime,
			).RowError(0, fmt.Errorf("row error"))

		mock.ExpectQuery(query).
			WillReturnRows(rows)

		ctx := context.Background()
		res, err := accessor.GetByLocation(ctx, location)

		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(res).To(gomega.BeNil())
	})
}
