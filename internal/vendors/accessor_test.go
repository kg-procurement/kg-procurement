package vendors

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"kg/procurement/internal/common/database"
	"log"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/benbjohnson/clock"
	"github.com/jmoiron/sqlx"
	"github.com/onsi/gomega"
)

func Test_newPostgresVendorAccessor(t *testing.T) {
	_ = newPostgresVendorAccessor(nil, nil)
}

func Test_postgresVendorAccessor_GetSomeStuff(t *testing.T) {
	t.Parallel()

	var (
		accessor *postgresVendorAccessor
		mock     sqlmock.Sqlmock
	)

	setup := func(t *testing.T) (*gomega.GomegaWithT, *sql.DB) {
		g := gomega.NewWithT(t)
		realClock := clock.New()
		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			log.Fatal("error initializing mock:", err)
		}
		sqlxDB := sqlx.NewDb(db, "sqlmock")

		accessor = newPostgresVendorAccessor(sqlxDB, realClock)
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

	t.Run("error on iterating row", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("FERRY").
			AddRow("VALE").
			RowError(1, fmt.Errorf("row error"))

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
		realClock := clock.New()
		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			log.Fatal("error initializing mock:", err)
		}
		sqlxDB := sqlx.NewDb(db, "sqlmock")

		accessor = newPostgresVendorAccessor(sqlxDB, realClock)
		mock = sqlMock

		return g, db
	}

	vendorFields := []string{
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

	dataQuery := `
		SELECT DISTINCT
			v.id,
			v.name,
			v.description,
			v.bp_id,
			v.bp_name,
			v.rating,
			v.area_group_id,
			v.area_group_name,
			v.sap_code,
			v.modified_date,
			v.modified_by,
			v.dt
		FROM vendor v
		ORDER BY v.dt DESC
		LIMIT $1
		OFFSET $2
	`

	countQuery := "SELECT COUNT(*) from vendor"

	fixedTime := time.Date(2024, time.September, 27, 12, 30, 0, 0, time.UTC)

	spec := GetAllVendorSpec{
		PaginationSpec: database.PaginationSpec{
			Order: "DESC",
			Limit: 10,
			Page:  1,
		},
	}

	t.Run("success", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(vendorFields).
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

		args := database.BuildPaginationArgs(spec.PaginationSpec)

		mock.ExpectQuery(dataQuery).
			WithArgs(args.Limit, args.Offset).
			WillReturnRows(rows)

		totalRows := sqlmock.NewRows([]string{"count"}).AddRow(1)

		mock.ExpectQuery(countQuery).WillReturnRows(totalRows)

		ctx := context.Background()
		res, err := accessor.GetAll(ctx, spec)

		vendorsExpectation := []Vendor{{
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
			ModifiedBy:    "1",
			Date:          fixedTime,
		}}

		expectation := &AccessorGetAllPaginationData{
			Vendors:  vendorsExpectation,
			Metadata: res.Metadata,
		}

		g.Expect(err).To(gomega.BeNil())
		g.Expect(res).To(gomega.Equal(expectation))
	})

	t.Run("success with orderby", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(vendorFields).
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

		customQuery := `
			SELECT DISTINCT
				v.id,
				v.name,
				v.description,
				v.bp_id,
				v.bp_name,
				v.rating,
				v.area_group_id,
				v.area_group_name,
				v.sap_code,
				v.modified_date,
				v.modified_by,
				v.dt
			FROM vendor v
			ORDER BY v.rating DESC
			LIMIT $1
			OFFSET $2
		`

		customSpec := GetAllVendorSpec{
			PaginationSpec: database.PaginationSpec{
				Order:   "DESC",
				Limit:   10,
				Page:    1,
				OrderBy: "rating",
			},
		}

		args := database.BuildPaginationArgs(customSpec.PaginationSpec)

		mock.ExpectQuery(customQuery).
			WithArgs(args.Limit, args.Offset).
			WillReturnRows(rows)

		totalRows := sqlmock.NewRows([]string{"count"}).AddRow(1)

		mock.ExpectQuery(countQuery).WillReturnRows(totalRows)

		ctx := context.Background()
		res, err := accessor.GetAll(ctx, customSpec)

		vendorsExpectation := []Vendor{{
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
			ModifiedBy:    "1",
			Date:          fixedTime,
		}}

		expectation := &AccessorGetAllPaginationData{
			Vendors:  vendorsExpectation,
			Metadata: res.Metadata,
		}

		g.Expect(err).To(gomega.BeNil())
		g.Expect(res).To(gomega.Equal(expectation))
	})

	t.Run("success with all pages fully filled", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(vendorFields).
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

		customSpec := GetAllVendorSpec{
			PaginationSpec: database.PaginationSpec{
				Order: "DESC",
				Limit: 1,
				Page:  1,
			},
		}

		mock.ExpectQuery(dataQuery).
			WithArgs(1, 0).
			WillReturnRows(rows)

		totalRows := sqlmock.NewRows([]string{"count"}).AddRow(1)

		mock.ExpectQuery(countQuery).WillReturnRows(totalRows)

		ctx := context.Background()
		res, err := accessor.GetAll(ctx, customSpec)

		vendorsExpectation := []Vendor{{
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
			ModifiedBy:    "1",
			Date:          fixedTime,
		}}

		expectation := &AccessorGetAllPaginationData{
			Vendors:  vendorsExpectation,
			Metadata: res.Metadata,
		}

		g.Expect(err).To(gomega.BeNil())
		g.Expect(res).To(gomega.Equal(expectation))
	})

	t.Run("success on empty result", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(vendorFields)

		args := database.BuildPaginationArgs(spec.PaginationSpec)

		mock.ExpectQuery(dataQuery).
			WithArgs(args.Limit, args.Offset).
			WillReturnRows(rows)

		totalRows := sqlmock.NewRows([]string{"count"}).AddRow(0)

		mock.ExpectQuery(countQuery).
			WillReturnRows(totalRows)

		ctx := context.Background()
		res, err := accessor.GetAll(ctx, spec)

		expectation := &AccessorGetAllPaginationData{
			Vendors:  []Vendor{},
			Metadata: res.Metadata,
		}

		g.Expect(err).To(gomega.BeNil())
		g.Expect(res).To(gomega.Equal(expectation))
	})

	t.Run("error on scanning vendor data rows", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(vendorFields).AddRow(
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

		args := database.BuildPaginationArgs(spec.PaginationSpec)

		mock.ExpectQuery(dataQuery).
			WithArgs(args.Limit, args.Offset).
			WillReturnRows(rows)

		totalRows := sqlmock.NewRows([]string{"count"}).AddRow(nil)

		mock.ExpectQuery(countQuery).WillReturnRows(totalRows)

		ctx := context.Background()
		res, err := accessor.GetAll(ctx, spec)

		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(res).To(gomega.BeNil())
	})

	t.Run("error on executing vendor data query", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(vendorFields).
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
			FROM vendor
			ORDER BY v.rating DESC
			LIMIT $1
			OFFSET $2
			`

		args := database.BuildPaginationArgs(spec.PaginationSpec)

		mock.ExpectQuery(wrongQuery).
			WithArgs(args.Limit, args.Offset).
			WillReturnRows(rows)

		totalRows := sqlmock.NewRows([]string{"count"}).AddRow(1)

		mock.ExpectQuery(countQuery).WillReturnRows(totalRows)

		ctx := context.Background()
		res, err := accessor.GetAll(ctx, spec)

		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(res).To(gomega.BeNil())
	})

	t.Run("error while iterating rows", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(vendorFields).
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

		args := database.BuildPaginationArgs(spec.PaginationSpec)

		mock.ExpectQuery(dataQuery).
			WithArgs(args.Limit, args.Offset).
			WillReturnRows(rows)

		totalRows := sqlmock.NewRows([]string{"count"}).AddRow(1)

		mock.ExpectQuery(countQuery).WillReturnRows(totalRows)

		ctx := context.Background()
		res, err := accessor.GetAll(ctx, spec)

		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(res).To(gomega.BeNil())
	})

	t.Run("error on scanning total entry row", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(vendorFields).
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

		args := database.BuildPaginationArgs(spec.PaginationSpec)

		mock.ExpectQuery(dataQuery).
			WithArgs(args.Limit, args.Offset).
			WillReturnRows(rows)

		totalRows := sqlmock.NewRows([]string{"count"}).RowError(1, fmt.Errorf("row error"))

		mock.ExpectQuery(countQuery).WillReturnRows(totalRows)

		ctx := context.Background()
		res, err := accessor.GetAll(ctx, spec)

		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(res).To(gomega.BeNil())
	})

}

func Test_postgresVendorAccessor_GetAll_WithLocationAndProduct(t *testing.T) {
	t.Parallel()

	var (
		accessor *postgresVendorAccessor
		mock     sqlmock.Sqlmock
	)

	setup := func(t *testing.T) (*gomega.GomegaWithT, *sql.DB) {
		g := gomega.NewWithT(t)
		realClock := clock.New()

		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			log.Fatal("error initializing mock:", err)
		}
		sqlxDB := sqlx.NewDb(db, "sqlmock")

		accessor = newPostgresVendorAccessor(sqlxDB, realClock)
		mock = sqlMock

		return g, db
	}

	fixedTime := time.Date(2024, time.September, 27, 12, 30, 0, 0, time.UTC)

	spec := GetAllVendorSpec{
		PaginationSpec: database.PaginationSpec{
			Order: "DESC",
			Limit: 10,
			Page:  1,
		},
		Location: "Indonesia",
		Product:  "test product",
	}

	vendorFields := []string{
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

	countQuery := "SELECT COUNT(*) from vendor"

	dataQuery := fmt.Sprintf(`
		SELECT DISTINCT
			v.id,
			v.name,
			v.description,
			v.bp_id,
			v.bp_name,
			v.rating,
			v.area_group_id,
			v.area_group_name,
			v.sap_code,
			v.modified_date,
			v.modified_by,
			v.dt
		FROM vendor v
		JOIN price pr ON pr.vendor_id = v.id
		JOIN product_vendor pv ON pv.id = pr.product_vendor_id
		JOIN product p ON p.id = pv.product_id
		WHERE v.area_group_name = $1 AND p.name iLIKE $2 AND p.name iLIKE $3
		ORDER BY v.dt %s
		LIMIT $4
		OFFSET $5
	`, spec.PaginationSpec.Order)

	productNameList := strings.Fields(spec.Product)

	t.Run("success with location and product", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(vendorFields).
			AddRow(
				"1",
				"name",
				"description",
				1,
				"bp_name",
				1,
				1,
				"Indonesia",
				"sap_code",
				fixedTime,
				1,
				fixedTime,
			)

		args := database.BuildPaginationArgs(spec.PaginationSpec)

		mock.ExpectQuery(dataQuery).
			WithArgs(
				spec.Location,
				"%"+productNameList[0]+"%",
				"%"+productNameList[1]+"%",
				args.Limit,
				args.Offset,
			).
			WillReturnRows(rows)

		totalRows := sqlmock.NewRows([]string{"count"}).AddRow(1)

		mock.ExpectQuery(countQuery).WillReturnRows(totalRows)

		ctx := context.Background()
		res, err := accessor.GetAll(ctx, spec)

		vendorsExpectation := []Vendor{{
			ID:            "1",
			Name:          "name",
			Description:   "description",
			BpID:          "1",
			BpName:        "bp_name",
			Rating:        1,
			AreaGroupID:   "1",
			AreaGroupName: "Indonesia",
			SapCode:       "sap_code",
			ModifiedDate:  fixedTime,
			ModifiedBy:    "1",
			Date:          fixedTime,
		}}

		expectation := &AccessorGetAllPaginationData{
			Vendors:  vendorsExpectation,
			Metadata: res.Metadata,
		}

		g.Expect(err).To(gomega.BeNil())
		g.Expect(res).To(gomega.Equal(expectation))
	})

	t.Run("error on executing query with location and product", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		mock.ExpectQuery(dataQuery).
			WithArgs(
				spec.Location,
				"%"+productNameList[0]+"%",
				"%"+productNameList[1]+"%",
				10,
				0,
			).
			WillReturnError(errors.New("some error"))

		ctx := context.Background()
		res, err := accessor.GetAll(ctx, spec)

		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(res).To(gomega.BeNil())
	})

	t.Run("success with only location", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		spec := GetAllVendorSpec{
			PaginationSpec: database.PaginationSpec{
				Order: "DESC",
				Limit: 10,
				Page:  1,
			},
			Location: "Indonesia",
		}

		dataQuery := `
			SELECT DISTINCT
				v.id,
				v.name,
				v.description,
				v.bp_id,
				v.bp_name,
				v.rating,
				v.area_group_id,
				v.area_group_name,
				v.sap_code,
				v.modified_date,
				v.modified_by,
				v.dt
			FROM vendor v
			WHERE v.area_group_name = $1
			ORDER BY v.dt DESC
			LIMIT $2
			OFFSET $3
		`

		rows := sqlmock.NewRows(vendorFields).
			AddRow(
				"1",
				"name",
				"description",
				1,
				"bp_name",
				1,
				1,
				"Indonesia",
				"sap_code",
				fixedTime,
				1,
				fixedTime,
			)

		args := database.BuildPaginationArgs(spec.PaginationSpec)

		mock.ExpectQuery(dataQuery).
			WithArgs(spec.Location, args.Limit, args.Offset).
			WillReturnRows(rows)

		totalRows := sqlmock.NewRows([]string{"count"}).AddRow(1)

		mock.ExpectQuery(countQuery).WillReturnRows(totalRows)

		ctx := context.Background()
		res, err := accessor.GetAll(ctx, spec)

		vendorsExpectation := []Vendor{{
			ID:            "1",
			Name:          "name",
			Description:   "description",
			BpID:          "1",
			BpName:        "bp_name",
			Rating:        1,
			AreaGroupID:   "1",
			AreaGroupName: "Indonesia",
			SapCode:       "sap_code",
			ModifiedDate:  fixedTime,
			ModifiedBy:    "1",
			Date:          fixedTime,
		}}

		expectation := &AccessorGetAllPaginationData{
			Vendors:  vendorsExpectation,
			Metadata: res.Metadata,
		}

		g.Expect(err).To(gomega.BeNil())
		g.Expect(res).To(gomega.Equal(expectation))
	})

	t.Run("success with only product", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		spec := GetAllVendorSpec{
			PaginationSpec: database.PaginationSpec{
				Order: "DESC",
				Limit: 10,
				Page:  1,
			},
			Product: "test product",
		}

		dataQuery := `
			SELECT DISTINCT
				v.id,
				v.name,
				v.description,
				v.bp_id,
				v.bp_name,
				v.rating,
				v.area_group_id,
				v.area_group_name,
				v.sap_code,
				v.modified_date,
				v.modified_by,
				v.dt
			FROM vendor v
			JOIN price pr ON pr.vendor_id = v.id
			JOIN product_vendor pv ON pv.id = pr.product_vendor_id
			JOIN product p ON p.id = pv.product_id
			WHERE p.name iLIKE $1 AND p.name iLIKE $2
			ORDER BY v.dt DESC
			LIMIT $3
			OFFSET $4
		`

		rows := sqlmock.NewRows(vendorFields).
			AddRow(
				"1",
				"name",
				"description",
				1,
				"bp_name",
				1,
				1,
				"Indonesia",
				"sap_code",
				fixedTime,
				1,
				fixedTime,
			)

		args := database.BuildPaginationArgs(spec.PaginationSpec)

		mock.ExpectQuery(dataQuery).
			WithArgs(
				"%"+productNameList[0]+"%",
				"%"+productNameList[1]+"%",
				args.Limit,
				args.Offset,
			).
			WillReturnRows(rows)

		totalRows := sqlmock.NewRows([]string{"count"}).AddRow(1)

		mock.ExpectQuery(countQuery).WillReturnRows(totalRows)

		ctx := context.Background()
		res, err := accessor.GetAll(ctx, spec)

		vendorsExpectation := []Vendor{{
			ID:            "1",
			Name:          "name",
			Description:   "description",
			BpID:          "1",
			BpName:        "bp_name",
			Rating:        1,
			AreaGroupID:   "1",
			AreaGroupName: "Indonesia",
			SapCode:       "sap_code",
			ModifiedDate:  fixedTime,
			ModifiedBy:    "1",
			Date:          fixedTime,
		}}

		expectation := &AccessorGetAllPaginationData{
			Vendors:  vendorsExpectation,
			Metadata: res.Metadata,
		}

		g.Expect(err).To(gomega.BeNil())
		g.Expect(res).To(gomega.Equal(expectation))
	})
}

func TestVendorAccessor_GetById(t *testing.T) {
	t.Parallel()

	var (
		accessor *postgresVendorAccessor
		mock     sqlmock.Sqlmock
	)

	setup := func(t *testing.T) (*gomega.GomegaWithT, *sql.DB) {
		g := gomega.NewWithT(t)
		realClock := clock.New()

		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			log.Fatal("error initializing mock:", err)
		}
		sqlxDB := sqlx.NewDb(db, "sqlmock")

		accessor = newPostgresVendorAccessor(sqlxDB, realClock)
		mock = sqlMock

		return g, db
	}

	vendorFields := []string{
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
	WHERE id = $1`

	fixedTime := time.Date(2024, time.September, 27, 12, 30, 0, 0, time.UTC)

	t.Run("success", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(vendorFields).
			AddRow(
				"ID",
				"name",
				"description",
				"BpID",
				"BpName",
				1,
				"AreaGroupId",
				"AreaGroupName",
				"SapCode",
				fixedTime,
				"ID",
				fixedTime,
			)

		mock.ExpectQuery(query).
			WithArgs("ID").
			WillReturnRows(rows)

		ctx := context.Background()
		res, err := accessor.GetById(ctx, "ID")

		expectation := &Vendor{
			ID:            "ID",
			Name:          "name",
			Description:   "description",
			BpID:          "BpID",
			BpName:        "BpName",
			Rating:        1,
			AreaGroupID:   "AreaGroupId",
			AreaGroupName: "AreaGroupName",
			SapCode:       "SapCode",
			ModifiedDate:  fixedTime,
			ModifiedBy:    "ID",
			Date:          fixedTime,
		}

		g.Expect(err).To(gomega.BeNil())
		g.Expect(res).To(gomega.Equal(expectation))

	})

	t.Run("error on row scan", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(vendorFields).RowError(1, fmt.Errorf("error"))

		mock.ExpectQuery(query).
			WithArgs("ID").
			WillReturnRows(rows)

		ctx := context.Background()
		res, err := accessor.GetById(ctx, "ID")

		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(res).To(gomega.BeNil())

	})
}

func TestVendorAccessor_UpdateDetail(t *testing.T) {
	t.Parallel()

	var (
		accessor  *postgresVendorAccessor
		mock      sqlmock.Sqlmock
		clockMock *clock.Mock
	)

	updatedFixedTime := time.Date(2024, time.September, 27, 12, 30, 0, 1, time.UTC)

	fixedTime := time.Date(2024, time.September, 27, 12, 30, 0, 0, time.UTC)

	setup := func(t *testing.T) (*gomega.GomegaWithT, *sql.DB) {
		g := gomega.NewWithT(t)

		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			log.Fatal("error initializing mock:", err)
		}
		sqlxDB := sqlx.NewDb(db, "sqlmock")

		clockMock = clock.NewMock()

		accessor = newPostgresVendorAccessor(sqlxDB, clockMock)

		mock = sqlMock

		return g, db
	}

	vendorFields := []string{
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

	query := `UPDATE vendor
		SET 
			name = $2,
			description = $3,
			bp_id = $4,
			bp_name = $5,
			rating = $6,
			area_group_id = $7,
			area_group_name = $8,
			sap_code = $9,
			modified_date = $10
		WHERE 
			id = $1
		RETURNING 
			id, 
			name, 
			description, 
			bp_id, 
			bp_name, 
			rating, 
			area_group_id, 
			area_group_name, 
			sap_code, 
			modified_date, 
			modified_by, 
			dt
	`

	t.Run("success", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		clockMock.Set(time.Now())

		now := clockMock.Now()

		rows := sqlmock.NewRows(vendorFields).
			AddRow(
				"ID",
				"updated",
				"updated",
				"updated",
				"updated",
				2,
				"updated",
				"updated",
				"updated",
				now,
				"updater",
				fixedTime,
			)

		ctx := context.Background()

		updatedVendor := &Vendor{
			ID:            "ID",
			Name:          "updated",
			Description:   "updated",
			BpID:          "updated",
			BpName:        "updated",
			Rating:        2,
			AreaGroupID:   "updated",
			AreaGroupName: "updated",
			SapCode:       "updated",
			ModifiedDate:  now,
			ModifiedBy:    "updater",
			Date:          fixedTime,
		}

		mock.ExpectQuery(query).
			WithArgs("ID",
				"updated",
				"updated",
				"updated",
				"updated",
				2,
				"updated",
				"updated",
				"updated",
				now).
			WillReturnRows(rows)

		res, err := accessor.UpdateDetail(ctx, *updatedVendor)

		g.Expect(err).To(gomega.BeNil())
		g.Expect(res).To(gomega.Equal(updatedVendor))
	})

	t.Run("error on row scan", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(vendorFields).
			AddRow(
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
			WithArgs("ID",
				"updated",
				"updated",
				"updated",
				"updated",
				"not int",
				"updated",
				"updated",
				"updated",
				updatedFixedTime,
				"updatedID",
				fixedTime).
			WillReturnRows(rows)

		ctx := context.Background()
		res, err := accessor.UpdateDetail(ctx, Vendor{})

		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(res).To(gomega.BeNil())

	})
}

func Test_GetAllLocations(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupVendorAccessorTestComponent(t, WithQueryMatcher(sqlmock.QueryMatcherRegexp))
		)
		defer c.db.Close()

		rows := sqlmock.NewRows([]string{"area_group_name"}).
			AddRow("Location1").
			AddRow("Location2")

		c.mock.ExpectQuery(getAllLocationsQuery).WillReturnRows(rows)

		results, err := c.accessor.GetAllLocations(ctx)

		c.g.Expect(err).To(gomega.BeNil())
		c.g.Expect(results).To(gomega.Equal([]string{"Location1", "Location2"}))
	})

	t.Run("error - query execution", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupVendorAccessorTestComponent(t, WithQueryMatcher(sqlmock.QueryMatcherRegexp))
		)
		defer c.db.Close()

		c.mock.ExpectQuery(getAllLocationsQuery).WillReturnError(sql.ErrConnDone)

		results, err := c.accessor.GetAllLocations(ctx)

		c.g.Expect(err).ToNot(gomega.BeNil())
		c.g.Expect(err).To(gomega.Equal(sql.ErrConnDone))
		c.g.Expect(results).To(gomega.BeNil())
	})

	t.Run("error - row scan", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupVendorAccessorTestComponent(t, WithQueryMatcher(sqlmock.QueryMatcherRegexp))
		)
		defer c.db.Close()

		rows := sqlmock.NewRows([]string{"area_group_name"}).
			AddRow(nil)

		c.mock.ExpectQuery(getAllLocationsQuery).WillReturnRows(rows)

		results, err := c.accessor.GetAllLocations(ctx)

		c.g.Expect(err).ToNot(gomega.BeNil())
		c.g.Expect(results).To(gomega.BeNil())
	})
}

func Test_BulkGetByIDs(t *testing.T) {
	t.Parallel()

	var (
		vendorColumns = []string{"id", "email"}
		vendors       = []Vendor{
			{
				ID:    "1",
				Email: "ferryganteng@outlook.com",
			},
			{
				ID:    "2",
				Email: "valenganteng@gmail.com",
			},
		}
	)

	t.Run("success", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupVendorAccessorTestComponent(t)
		)
		defer c.db.Close()

		rows := sqlmock.NewRows(vendorColumns).
			AddRow("1", "ferryganteng@outlook.com").
			AddRow("2", "valenganteng@gmail.com")

		vendorIDs := []string{"1", "2"}
		query, _, _ := sqlx.In(getBulkByID, vendorIDs)
		c.mock.ExpectQuery(query).
			WithArgs("1", "2").
			WillReturnRows(rows)

		results, err := c.accessor.BulkGetByIDs(ctx, vendorIDs)

		c.g.Expect(err).To(gomega.BeNil())
		c.g.Expect(results).To(gomega.Equal(vendors))
	})

	t.Run("error on rows", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupVendorAccessorTestComponent(t)
		)
		defer c.db.Close()

		rows := sqlmock.NewRows(vendorColumns).
			AddRow("1", "ferryganteng@outlook.com").
			AddRow("2", "valenganteng@gmail.com").
			RowError(1, errors.New("oh no"))

		vendorIDs := []string{"1", "2"}
		query, _, _ := sqlx.In(getBulkByID, vendorIDs)
		c.mock.ExpectQuery(query).
			WithArgs("1", "2").
			WillReturnRows(rows)

		results, err := c.accessor.BulkGetByIDs(ctx, vendorIDs)

		c.g.Expect(err).ToNot(gomega.BeNil())
		c.g.Expect(results).To(gomega.BeNil())
	})

	t.Run("error - query execution", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupVendorAccessorTestComponent(t)
		)
		defer c.db.Close()

		vendorIDs := []string{"1", "2"}
		query, _, _ := sqlx.In(getBulkByID, vendorIDs)
		c.mock.ExpectQuery(query).
			WillReturnError(sql.ErrConnDone)

		results, err := c.accessor.BulkGetByIDs(ctx, vendorIDs)

		c.g.Expect(err).ToNot(gomega.BeNil())
		c.g.Expect(results).To(gomega.BeNil())
	})

	t.Run("error - row scan", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupVendorAccessorTestComponent(t)
		)
		defer c.db.Close()

		rows := sqlmock.NewRows(vendorColumns).
			AddRow(nil, nil)

		c.mock.ExpectQuery(getBulkByID).
			WillReturnRows(rows)

		results, err := c.accessor.BulkGetByIDs(ctx, []string{"1"})
		c.g.Expect(err).ToNot(gomega.BeNil())
		c.g.Expect(results).To(gomega.BeNil())
	})

	t.Run("error - empty slice", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupVendorAccessorTestComponent(t)
		)
		defer c.db.Close()

		rows := sqlmock.NewRows(vendorColumns).
			AddRow("1", "ferryganteng@outlook.com").
			AddRow("2", "valenganteng@gmail.com")

		c.mock.ExpectQuery(getBulkByID).
			WillReturnRows(rows)

		results, err := c.accessor.BulkGetByIDs(ctx, []string{})
		c.g.Expect(err).ToNot(gomega.BeNil())
		c.g.Expect(results).To(gomega.BeNil())
	})
}

func Test_writeProductVendor(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var (
			ctx    = context.Background()
			c      = setupVendorAccessorTestComponent(t, WithQueryMatcher(sqlmock.QueryMatcherRegexp))
			vendor = Vendor{ID: "123"}
		)

		transformedQuery, args, _ := sqlx.Named(insertVendor, vendor)
		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		c.mock.ExpectExec(regexp.QuoteMeta(transformedQuery)).WithArgs(
			driverArgs...,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err := c.accessor.writeVendor(ctx, vendor)
		c.g.Expect(err).Should(gomega.BeNil())
	})

	t.Run("error", func(t *testing.T) {
		var (
			ctx    = context.Background()
			c      = setupVendorAccessorTestComponent(t)
			vendor = Vendor{ID: "123"}
		)

		transformedQuery, args, _ := sqlx.Named(insertVendor, vendor)
		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		c.mock.ExpectExec(regexp.QuoteMeta(transformedQuery)).WithArgs(
			driverArgs...,
		).WillReturnError(errors.New("error"))

		err := c.accessor.writeVendor(ctx, vendor)
		c.g.Expect(err).ShouldNot(gomega.BeNil())
	})
}

func Test_getAllVendorIdByProductName(t *testing.T) {
	t.Parallel()

	var (
		productName   = "book"
		vendorColumns = []string{"id", "email"}
		vendors       = []Vendor{
			{
				ID:    "1",
				Email: "valerian@outlook.com",
			},
			{
				ID:    "2",
				Email: "salim@outlook.com",
			},
		}
	)

	t.Run("success", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupVendorAccessorTestComponent(t)
		)
		defer c.db.Close()

		// expected := []string{"2550", "2580"}

		transformedQuery, args, _ := sqlx.Named(getBulkByProductName, map[string]interface{}{
			"product_name": productName,
		})

		rows := sqlmock.NewRows(vendorColumns).
			AddRow("1", "valerian@outlook.com").
			AddRow("2", "salim@outlook.com")

		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		c.mock.ExpectQuery(transformedQuery).
			WithArgs(driverArgs...).
			WillReturnRows(rows)

		results, err := c.accessor.BulkGetByProductName(ctx, productName)

		c.g.Expect(err).To(gomega.BeNil())
		c.g.Expect(results).To(gomega.Equal(vendors))
	})

	t.Run("success returning empty array", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupVendorAccessorTestComponent(t)
		)
		defer c.db.Close()

		expected := []Vendor{}

		rows := sqlmock.NewRows(vendorColumns)

		transformedQuery, args, _ := sqlx.Named(getBulkByProductName, map[string]interface{}{
			"product_name": productName,
		})

		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		c.mock.ExpectQuery(transformedQuery).
			WithArgs(driverArgs...).
			WillReturnRows(rows)

		results, err := c.accessor.BulkGetByProductName(ctx, productName)

		c.g.Expect(err).To(gomega.BeNil())
		c.g.Expect(results).To(gomega.Equal(expected))
	})

	t.Run("errors on db query", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupVendorAccessorTestComponent(t)
		)
		defer c.db.Close()

		transformedQuery, args, _ := sqlx.Named(getBulkByProductName, map[string]interface{}{
			"product_name": productName,
		})

		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		c.mock.ExpectQuery(transformedQuery).
			WithArgs(driverArgs...).
			WillReturnError(errors.New("error"))

		results, err := c.accessor.BulkGetByProductName(ctx, productName)

		c.g.Expect(err).ToNot(gomega.BeNil())
		c.g.Expect(results).To(gomega.BeNil())
	})

	t.Run("errors while iterating rows", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupVendorAccessorTestComponent(t)
		)
		defer c.db.Close()

		rows := sqlmock.NewRows(vendorColumns).
			AddRow("1", "valerian@outlook.com").
			AddRow("2", "salim@outlook.com").
			RowError(1, errors.New("error"))

		transformedQuery, args, _ := sqlx.Named(getBulkByProductName, map[string]interface{}{
			"product_name": productName,
		})

		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		c.mock.ExpectQuery(transformedQuery).
			WithArgs(driverArgs...).
			WillReturnRows(rows)

		results, err := c.accessor.BulkGetByProductName(ctx, productName)

		c.g.Expect(err).ToNot(gomega.BeNil())
		c.g.Expect(results).To(gomega.BeNil())
	})

	t.Run("errors on scanning rows", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupVendorAccessorTestComponent(t)
		)
		defer c.db.Close()

		rows := sqlmock.NewRows(vendorColumns).
			AddRow(nil, nil)

		transformedQuery, args, _ := sqlx.Named(getBulkByProductName, map[string]interface{}{
			"product_name": productName,
		})

		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		c.mock.ExpectQuery(transformedQuery).
			WithArgs(driverArgs...).
			WillReturnRows(rows)

		results, err := c.accessor.BulkGetByProductName(ctx, productName)

		c.g.Expect(err).ToNot(gomega.BeNil())
		c.g.Expect(results).To(gomega.BeNil())
	})
}

func Test_createVendorEvaluation(t *testing.T) {
	t.Parallel()

	var (
		vendorEvaluation = VendorEvaluation{
			ID:                               "Vp5XxWumASFvxD3",
			VendorID:                         "1",
			KesesuaianProduk:                 1,
			KualitasProduk:                   1,
			KetepatanWaktuPengiriman:         1,
			KompetitifitasHarga:              1,
			ResponsivitasKemampuanKomunikasi: 1,
			KemampuanDalamMenanganiMasalah:   1,
			KelengkapanBarang:                1,
			Harga:                            1,
			TermOfPayment:                    1,
			Reputasi:                         1,
			KetersediaanBarang:               1,
			KualitasLayananAfterServices:     1,
		}
	)

	t.Run("success", func(t *testing.T) {
		var (
			c   = setupVendorAccessorTestComponent(t)
			ctx = context.Background()
		)

		transformedQuery, args, _ := sqlx.Named(createEvaluationQuery, vendorEvaluation)

		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		c.mock.ExpectExec(transformedQuery).
			WithArgs(driverArgs...).WillReturnResult(sqlmock.NewResult(1, 1))

		_, err := c.accessor.CreateEvaluation(ctx, &vendorEvaluation)

		c.g.Expect(err).To(gomega.BeNil())
	})

	t.Run("error while doing db query", func(t *testing.T) {
		var (
			c   = setupVendorAccessorTestComponent(t)
			ctx = context.Background()
		)

		transformedQuery, args, _ := sqlx.Named(createEvaluationQuery, vendorEvaluation)

		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		c.mock.ExpectExec(transformedQuery).
			WithArgs(driverArgs...).WillReturnError(sql.ErrConnDone)

		_, err := c.accessor.CreateEvaluation(ctx, &vendorEvaluation)

		c.g.Expect(err).ToNot(gomega.BeNil())
	})

}

func Test_Close(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		c := setupVendorAccessorTestComponent(t)
		c.mock.ExpectClose()
		err := c.accessor.Close()
		c.g.Expect(err).Should(gomega.BeNil())
	})

	t.Run("error", func(t *testing.T) {
		c := setupVendorAccessorTestComponent(t)
		c.mock.ExpectClose().WillReturnError(errors.New("error"))
		err := c.accessor.Close()
		c.g.Expect(err).ShouldNot(gomega.BeNil())
	})
}

type vendorAccessorTestComponent struct {
	g        *gomega.WithT
	mock     sqlmock.Sqlmock
	db       *sql.DB
	accessor *postgresVendorAccessor
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

func setupVendorAccessorTestComponent(t *testing.T, opts ...Option) vendorAccessorTestComponent {
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

	return vendorAccessorTestComponent{
		g:        g,
		mock:     sqlMock,
		db:       db,
		accessor: newPostgresVendorAccessor(sqlxDB, clockMock),
		cmock:    clockMock,
	}
}
