package vendors

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"kg/procurement/internal/common/database"
	"log"
	"strings"
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
		ORDER BY rating DESC
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
			ModifiedBy:    1,
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
			ModifiedBy:    1,
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
			ORDER BY rating DESC
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

		mock.ExpectQuery(query).
			WithArgs(location).
			WillReturnError(errors.New("some error"))

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
			WithArgs(location).
			WillReturnRows(rows)

		ctx := context.Background()
		res, err := accessor.GetByLocation(ctx, location)

		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(res).To(gomega.BeNil())
	})
}

func TestVendorAccessor_GetByProductDescription(t *testing.T) {
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
			WHERE description LIKE $1 AND description LIKE $2`

	fixedTime := time.Date(2024, time.September, 27, 12, 30, 0, 0, time.UTC)
	product := "test product"
	productDescription := strings.Fields(product)

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
				"group_name",
				"sap_code",
				fixedTime,
				1,
				fixedTime,
			)

		mock.ExpectQuery(query).
			WithArgs("%"+productDescription[0]+"%", "%"+productDescription[1]+"%").
			WillReturnRows(rows)

		ctx := context.Background()
		res, err := accessor.GetByProductDescription(ctx, productDescription)

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
			WithArgs("%"+productDescription[0]+"%", "%"+productDescription[1]+"%").
			WillReturnRows(rows)

		ctx := context.Background()
		res, err := accessor.GetByProductDescription(ctx, productDescription)
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
			WithArgs("%"+productDescription[0]+"%", "%"+productDescription[1]+"%").
			WillReturnRows(rows)

		ctx := context.Background()
		res, err := accessor.GetByProductDescription(ctx, productDescription)

		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(res).To(gomega.BeNil())
	})

	t.Run("error on executing query", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		mock.ExpectQuery(query).
			WithArgs("%"+productDescription[0]+"%", "%"+productDescription[1]+"%").
			WillReturnError(errors.New("some error"))

		ctx := context.Background()
		res, err := accessor.GetByProductDescription(ctx, productDescription)

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
				"group_name",
				"sap_code",
				fixedTime,
				1,
				fixedTime,
			).RowError(0, fmt.Errorf("row error"))

		mock.ExpectQuery(query).
			WithArgs("%"+productDescription[0]+"%", "%"+productDescription[1]+"%").
			WillReturnRows(rows)

		ctx := context.Background()
		res, err := accessor.GetByProductDescription(ctx, productDescription)

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
		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			log.Fatal("error initializing mock:", err)
		}

		accessor = newPostgresVendorAccessor(db)
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
		JOIN product_vendor pv ON pv.vendor_id = v.id
		JOIN product p ON p.id = pv.product_id
		WHERE area_group_name = $1 AND p.name LIKE $2 AND p.name LIKE $3
		ORDER BY rating %s
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
			ModifiedBy:    1,
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
			WHERE area_group_name = $1
			ORDER BY rating DESC
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
			ModifiedBy:    1,
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
			JOIN product_vendor pv ON pv.vendor_id = v.id
			JOIN product p ON p.id = pv.product_id
			WHERE p.name LIKE $1 AND p.name LIKE $2
			ORDER BY rating DESC
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
			ModifiedBy:    1,
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
