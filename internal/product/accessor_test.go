package product

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

func Test_newPostgresProductAccessor(t *testing.T) {
	_ = newPostgresProductAccessor(nil, nil)
}

func Test_UpdateProduct(t *testing.T) {
	t.Parallel()

	productFields := []string{
		"id",
		"product_category_id",
		"uom_id",
		"income_tax_id",
		"product_type_id",
		"name",
		"description",
		"modified_date",
		"modified_by",
	}

	t.Run("success", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupProductAccessorTestComponent(t)
		)

		setup := func(t *testing.T) (*gomega.GomegaWithT, *sqlx.DB) {
			g := gomega.NewWithT(t)
			c.cmock.Set(time.Now())

			db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				log.Fatal("error initializing mock:", err)
			}
			sqlxDB := sqlx.NewDb(db, "sqlmock")

			return g, sqlxDB
		}

		t.Run("When updating product successfully, should return no error", func(t *testing.T) {
			_, db := setup(t)
			defer db.Close()

			now := c.cmock.Now()
			var (
				expectedResult = sqlmock.NewRows(productFields).
					AddRow(
						"inv1",
						"category_id_updated",
						"uom_id_updated",
						"income_tax_id_updated",
						"product_type_id_updated",
						"Product A Updated",
						"Updated Description",
						now,
						"someone",
					)
			)
			defer c.db.Close()

			updatedProduct := Product{
				ID:                "inv1",
				ProductCategoryID: "category_id_updated",
				UOMID:             "uom_id_updated",
				IncomeTaxID:       "income_tax_id_updated",
				ProductTypeID:     "product_type_id_updated",
				Name:              "Product A Updated",
				Description:       "Updated Description",
				ModifiedBy:        "someone",
				ModifiedDate:      now,
			}

			c.mock.ExpectQuery(updateProduct).
				WithArgs(
					updatedProduct.ID,
					updatedProduct.ProductCategoryID,
					updatedProduct.UOMID,
					updatedProduct.IncomeTaxID,
					updatedProduct.ProductTypeID,
					updatedProduct.Name,
					updatedProduct.Description,
					now,
				).WillReturnRows(expectedResult)

			res, err := c.accessor.UpdateProduct(ctx, updatedProduct)

			c.g.Expect(err).To(gomega.BeNil())
			c.g.Expect(res).To(gomega.Equal(updatedProduct))
		})

		t.Run("error on row scan", func(t *testing.T) {
			var (
				ctx            = context.Background()
				c              = setupProductAccessorTestComponent(t)
				expectedResult = sqlmock.NewRows(productFields).
						AddRow(
						"nil",
						"nil",
						"nil",
						"nil",
						"nil",
						"nil",
						"nil",
						"nil",
						"nil",
					)
			)
			defer c.db.Close()

			c.mock.ExpectQuery(updateProduct).
				WithArgs(
					"inv1",
					"category_id_updated",
					"uom_id_updated",
					"income_tax_id_updated",
					"product_type_id_updated",
					"Product A Updated",
					"Updated Description",
					"not a time",
				).WillReturnRows(expectedResult)

			res, err := c.accessor.UpdateProduct(ctx, Product{})

			c.g.Expect(err).ToNot(gomega.BeNil())
			c.g.Expect(res).To(gomega.Equal(Product{}))
		})
	})
}

func Test_UpdatePrice(t *testing.T) {
	t.Parallel()

	priceFields := []string{
		"id",
		"purchasing_org_id",
		"vendor_id",
		"product_vendor_id",
		"quantity_min",
		"quantity_max",
		"quantity_uom_id",
		"lead_time_min",
		"lead_time_max",
		"currency_id",
		"price",
		"price_quantity",
		"price_uom_id",
		"valid_from",
		"valid_to",
		"valid_pattern_id",
		"area_group_id",
		"reference_number",
		"reference_date",
		"document_type_id",
		"document_id",
		"item_id",
		"term_of_payment_id",
		"invocation_order",
		"modified_date",
		"modified_by",
	}

	query := `UPDATE price
        SET 
            purchasing_org_id = $2,
            vendor_id = $3,
            product_vendor_id = $4,
            quantity_min = $5,
            quantity_max = $6,
            quantity_uom_id = $7,
            lead_time_min = $8,
            lead_time_max = $9,
            currency_id = $10,
            price = $11,
            price_quantity = $12,
            price_uom_id = $13,
            valid_from = $14,
            valid_to = $15,
            valid_pattern_id = $16,
            area_group_id = $17,
            reference_number = $18,
            reference_date = $19,
            document_type_id = $20,
            document_id = $21,
            item_id = $22,
            term_of_payment_id = $23,
            invocation_order = $24,
            modified_date = $25
        WHERE 
            id = $1
        RETURNING 
            id,
            purchasing_org_id,
            vendor_id,
            product_vendor_id,
            quantity_min,
            quantity_max,
            quantity_uom_id,
            lead_time_min,
            lead_time_max,
            currency_id,
            price,
            price_quantity,
            price_uom_id,
            valid_from,
            valid_to,
            valid_pattern_id,
            area_group_id,
            reference_number,
            reference_date,
            document_type_id,
            document_id,
            item_id,
            term_of_payment_id,
            invocation_order,
            modified_date,
            modified_by
    `
	fixedTime := time.Now()
	updatedFixedTime := time.Now().Add(1 * time.Hour)

	t.Run("success", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupProductAccessorTestComponent(t)
		)
		c.cmock.Set(time.Now())
		now := c.cmock.Now()
		var (
			expectedResult = sqlmock.NewRows(priceFields).
				AddRow(
					"ID",
					"org_id_updated",
					"vendor_id_updated",
					"product_vendor_id_updated",
					10,
					100,
					"quantity_uom_id_updated",
					1,
					2,
					"currency_id_updated",
					99.99,
					1,
					"price_uom_id_updated",
					fixedTime,
					updatedFixedTime,
					"valid_pattern_id_updated",
					"area_group_id_updated",
					"reference_number_updated",
					fixedTime,
					"document_type_id_updated",
					"document_id_updated",
					"item_id_updated",
					"term_of_payment_id_updated",
					1,
					now,
					"modified_by_updated",
				)
		)
		defer c.db.Close()

		updatedPrice := Price{
			ID:              "ID",
			PurchasingOrgID: "org_id_updated",
			VendorID:        "vendor_id_updated",
			ProductVendorID: "product_vendor_id_updated",
			QuantityMin:     10,
			QuantityMax:     100,
			QuantityUOMID:   "quantity_uom_id_updated",
			LeadTimeMin:     1,
			LeadTimeMax:     2,
			CurrencyID:      "currency_id_updated",
			Price:           99.99,
			PriceQuantity:   1,
			PriceUOMID:      "price_uom_id_updated",
			ValidFrom:       fixedTime,
			ValidTo:         updatedFixedTime,
			ValidPatternID:  "valid_pattern_id_updated",
			AreaGroupID:     "area_group_id_updated",
			ReferenceNumber: "reference_number_updated",
			ReferenceDate:   fixedTime,
			DocumentTypeID:  "document_type_id_updated",
			DocumentID:      "document_id_updated",
			ItemID:          "item_id_updated",
			TermOfPaymentID: "term_of_payment_id_updated",
			InvocationOrder: 1,
			ModifiedDate:    now,
			ModifiedBy:      "modified_by_updated",
		}

		c.mock.ExpectQuery(query).
			WithArgs(
				"ID",
				updatedPrice.PurchasingOrgID,
				updatedPrice.VendorID,
				updatedPrice.ProductVendorID,
				updatedPrice.QuantityMin,
				updatedPrice.QuantityMax,
				updatedPrice.QuantityUOMID,
				updatedPrice.LeadTimeMin,
				updatedPrice.LeadTimeMax,
				updatedPrice.CurrencyID,
				updatedPrice.Price,
				updatedPrice.PriceQuantity,
				updatedPrice.PriceUOMID,
				updatedPrice.ValidFrom,
				updatedPrice.ValidTo,
				updatedPrice.ValidPatternID,
				updatedPrice.AreaGroupID,
				updatedPrice.ReferenceNumber,
				updatedPrice.ReferenceDate,
				updatedPrice.DocumentTypeID,
				updatedPrice.DocumentID,
				updatedPrice.ItemID,
				updatedPrice.TermOfPaymentID,
				updatedPrice.InvocationOrder,
				now,
			).WillReturnRows(expectedResult)

		res, err := c.accessor.UpdatePrice(ctx, updatedPrice)

		c.g.Expect(err).To(gomega.BeNil())
		c.g.Expect(res).To(gomega.Equal(updatedPrice))
	})

	t.Run("error on row scan", func(t *testing.T) {
		var (
			ctx            = context.Background()
			c              = setupProductAccessorTestComponent(t)
			expectedResult = sqlmock.NewRows(priceFields).
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
					nil,
					nil,
				)
		)
		defer c.db.Close()

		c.mock.ExpectQuery(query).
			WithArgs(
				"ID",
				"org_id_updated",
				"vendor_id_updated",
				"product_vendor_id_updated",
				10,
				100,
				"quantity_uom_id_updated",
				"not int",
				"not int",
				"currency_id_updated",
				99.99,
				1,
				"price_uom_id_updated",
				fixedTime,
				updatedFixedTime,
				"valid_pattern_id_updated",
				"area_group_id_updated",
				"reference_number_updated",
				fixedTime,
				"document_type_id_updated",
				"document_id_updated",
				"item_id_updated",
				"term_of_payment_id_updated",
				1,
				updatedFixedTime,
			).WillReturnRows(expectedResult)

		res, err := c.accessor.UpdatePrice(ctx, Price{})

		c.g.Expect(err).ToNot(gomega.BeNil())
		c.g.Expect(res).To(gomega.Equal(Price{}))
	})
}

func Test_GetProductVendorsByVendor(t *testing.T) {
	t.Parallel()

	var (
		now      = time.Now()
		vendorID = "1234"
		spec     = GetProductVendorByVendorSpec{
			PaginationSpec: database.PaginationSpec{
				Limit: 10,
				Page:  1,
			},
		}
		args           = database.BuildPaginationArgs(spec.PaginationSpec)
		productColumns = []string{"id", "product_id", "code", "name", "income_tax_id", "income_tax_name", "income_tax_percentage", "description", "uom_id", "sap_code", "modified_date", "modified_by"}
		productVendors = []ProductVendor{
			{
				ID:           "1111",
				Name:         "Mixer",
				ModifiedDate: now,
			},
			{
				ID:           "2222",
				Name:         "Rice Cooker",
				ModifiedDate: now,
			},
		}
	)

	t.Run("success", func(t *testing.T) {
		var (
			ctx            = context.Background()
			c              = setupProductAccessorTestComponent(t)
			expectedResult = sqlmock.NewRows(productColumns).
					AddRow(productVendors[0].ID, "", "", productVendors[0].Name, "", "", "", "", "", "", productVendors[0].ModifiedDate, "").
					AddRow(productVendors[1].ID, "", "", productVendors[1].Name, "", "", "", "", "", "", productVendors[1].ModifiedDate, "")
		)
		defer c.db.Close()

		query := getProductVendorsByVendorQuery + " LIMIT $2 OFFSET $3"
		c.mock.ExpectQuery(query).
			WithArgs(vendorID, args.Limit, args.Offset).
			WillReturnRows(expectedResult)

		totalRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
		countQuery := `SELECT COUNT(*)
			FROM product_vendor pv
			JOIN price pr ON pr.product_vendor_id = pv.id
			WHERE pr.vendor_id = $1
		`
		c.mock.ExpectQuery(countQuery).
			WithArgs(vendorID).
			WillReturnRows(totalRows)

		expect := &AccessorGetProductVendorsByVendorPaginationData{
			ProductVendors: productVendors,
			Metadata: database.PaginationMetadata{
				TotalPage:    1,
				CurrentPage:  1,
				TotalEntries: 2,
			},
		}

		res, err := c.accessor.GetProductVendorsByVendor(ctx, vendorID, spec)
		c.g.Expect(err).To(gomega.BeNil())
		c.g.Expect(res).To(gomega.BeComparableTo(expect))
	})

	t.Run("success with filter by name", func(t *testing.T) {
		var (
			ctx            = context.Background()
			c              = setupProductAccessorTestComponent(t)
			expectedResult = sqlmock.NewRows(productColumns).
					AddRow(productVendors[1].ID, "", "", productVendors[1].Name, "", "", "", "", "", "", productVendors[1].ModifiedDate, "")
			customSpec = GetProductVendorByVendorSpec{
				Name:           "Rice Cooker",
				PaginationSpec: spec.PaginationSpec,
			}
		)
		defer c.db.Close()

		productNameList := strings.Fields(customSpec.Name)
		c.mock.ExpectQuery(getProductVendorsByVendorQuery+" AND p.name iLIKE $2 AND p.name iLIKE $3"+" LIMIT $4 OFFSET $5").
			WithArgs(vendorID, "%"+productNameList[0]+"%", "%"+productNameList[1]+"%", args.Limit, args.Offset).
			WillReturnRows(expectedResult)

		totalRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		countQuery := `SELECT COUNT(*)
			FROM product_vendor pv
			JOIN price pr ON pr.product_vendor_id = pv.id
			WHERE pr.vendor_id = $1
		`

		c.mock.ExpectQuery(countQuery).
			WithArgs(vendorID).
			WillReturnRows(totalRows)

		expect := &AccessorGetProductVendorsByVendorPaginationData{
			ProductVendors: productVendors[1:],
			Metadata: database.PaginationMetadata{
				TotalPage:    1,
				CurrentPage:  1,
				TotalEntries: 1,
			},
		}

		res, err := c.accessor.GetProductVendorsByVendor(ctx, vendorID, customSpec)
		c.g.Expect(err).To(gomega.BeNil())
		c.g.Expect(res).To(gomega.BeComparableTo(expect))
	})

	t.Run("success with order by", func(t *testing.T) {
		var (
			ctx            = context.Background()
			c              = setupProductAccessorTestComponent(t)
			expectedResult = sqlmock.NewRows(productColumns).
					AddRow(productVendors[0].ID, "", "", productVendors[0].Name, "", "", "", "", "", "", productVendors[0].ModifiedDate, "").
					AddRow(productVendors[1].ID, "", "", productVendors[1].Name, "", "", "", "", "", "", productVendors[1].ModifiedDate, "")
			customSpec = GetProductVendorByVendorSpec{
				PaginationSpec: database.PaginationSpec{
					OrderBy: "name",
					Limit:   10,
					Page:    1,
				},
			}
		)
		defer c.db.Close()

		c.mock.ExpectQuery(getProductVendorsByVendorQuery+" ORDER BY name ASC LIMIT $2 OFFSET $3").
			WithArgs(vendorID, args.Limit, args.Offset).
			WillReturnRows(expectedResult)

		totalRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
		countQuery := `SELECT COUNT(*)
			FROM product_vendor pv
			JOIN price pr ON pr.product_vendor_id = pv.id
			WHERE pr.vendor_id = $1
		`

		c.mock.ExpectQuery(countQuery).
			WithArgs(vendorID).
			WillReturnRows(totalRows)

		expect := &AccessorGetProductVendorsByVendorPaginationData{
			ProductVendors: productVendors,
			Metadata: database.PaginationMetadata{
				TotalPage:    1,
				CurrentPage:  1,
				TotalEntries: 2,
			},
		}

		res, err := c.accessor.GetProductVendorsByVendor(ctx, vendorID, customSpec)
		c.g.Expect(err).To(gomega.BeNil())
		c.g.Expect(res).To(gomega.BeComparableTo(expect))
	})

	t.Run("success with order by and filter by name", func(t *testing.T) {
		var (
			ctx            = context.Background()
			c              = setupProductAccessorTestComponent(t)
			expectedResult = sqlmock.NewRows(productColumns).
					AddRow(productVendors[1].ID, "", "", productVendors[1].Name, "", "", "", "", "", "", productVendors[1].ModifiedDate, "")
			customSpec = GetProductVendorByVendorSpec{
				Name: "Rice Cooker",
				PaginationSpec: database.PaginationSpec{
					OrderBy: "name",
					Limit:   10,
					Page:    1,
				},
			}
		)
		defer c.db.Close()

		productNameList := strings.Fields(customSpec.Name)
		c.mock.ExpectQuery(getProductVendorsByVendorQuery+
			" AND p.name iLIKE $2 AND p.name iLIKE $3 ORDER BY name ASC LIMIT $4 OFFSET $5").
			WithArgs(vendorID, "%"+productNameList[0]+"%", "%"+productNameList[1]+"%", args.Limit, args.Offset).
			WillReturnRows(expectedResult)

		totalRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		countQuery := `SELECT COUNT(*)
			FROM product_vendor pv
			JOIN price pr ON pr.product_vendor_id = pv.id
			WHERE pr.vendor_id = $1
		`

		c.mock.ExpectQuery(countQuery).
			WithArgs(vendorID).
			WillReturnRows(totalRows)

		expect := &AccessorGetProductVendorsByVendorPaginationData{
			ProductVendors: productVendors[1:],
			Metadata: database.PaginationMetadata{
				TotalPage:    1,
				CurrentPage:  1,
				TotalEntries: 1,
			},
		}

		res, err := c.accessor.GetProductVendorsByVendor(ctx, vendorID, customSpec)
		c.g.Expect(err).To(gomega.BeNil())
		c.g.Expect(res).To(gomega.BeComparableTo(expect))
	})

	t.Run("error on query execution", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupProductAccessorTestComponent(t)
		)
		defer c.db.Close()

		c.mock.ExpectQuery(getProductVendorsByVendorQuery+" LIMIT $1 OFFSET $2").
			WithArgs(vendorID, args.Limit, args.Offset).
			WillReturnError(errors.New("error"))

		res, err := c.accessor.GetProductVendorsByVendor(ctx, vendorID, spec)

		c.g.Expect(err).ShouldNot(gomega.BeNil())
		c.g.Expect(res).To(gomega.BeNil())
	})

	t.Run("error on scanning data row", func(t *testing.T) {
		var (
			ctx            = context.Background()
			c              = setupProductAccessorTestComponent(t)
			expectedResult = sqlmock.NewRows(productColumns).
					AddRow(productVendors[0].ID, "", "", productVendors[0].Name, "", "", "", "", "", "", productVendors[0].ModifiedDate, "").
					AddRow(productVendors[1].ID, "", "", productVendors[1].Name, "", "", "", "", "", "", productVendors[1].ModifiedDate, "").
					AddRow(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
		)
		defer c.db.Close()

		c.mock.ExpectQuery(getProductVendorsByVendorQuery+" LIMIT $1 OFFSET $2").
			WithArgs(vendorID, args.Limit, args.Offset).
			WillReturnRows(expectedResult)

		res, err := c.accessor.GetProductVendorsByVendor(ctx, vendorID, spec)

		c.g.Expect(err).ShouldNot(gomega.BeNil())
		c.g.Expect(res).To(gomega.BeNil())
	})

	t.Run("error while iterating rows", func(t *testing.T) {
		var (
			ctx            = context.Background()
			c              = setupProductAccessorTestComponent(t)
			expectedResult = sqlmock.NewRows(productColumns).
					AddRow(productVendors[0].ID, "", "", productVendors[0].Name, "", "", "", "", "", "", productVendors[0].ModifiedDate, "").
					AddRow(productVendors[1].ID, "", "", productVendors[1].Name, "", "", "", "", "", "", productVendors[1].ModifiedDate, "").
					RowError(0, fmt.Errorf("some error"))
		)
		defer c.db.Close()

		c.mock.ExpectQuery(getProductVendorsByVendorQuery+" LIMIT $1 OFFSET $2").
			WithArgs(vendorID, args.Limit, args.Offset).
			WillReturnRows(expectedResult)

		res, err := c.accessor.GetProductVendorsByVendor(ctx, vendorID, spec)

		c.g.Expect(err).ToNot(gomega.BeNil())
		c.g.Expect(res).To(gomega.BeNil())
	})

	t.Run("error while querying count", func(t *testing.T) {
		var (
			ctx            = context.Background()
			c              = setupProductAccessorTestComponent(t)
			expectedResult = sqlmock.NewRows(productColumns).
					AddRow(productVendors[0].ID, "", "", productVendors[0].Name, "", "", "", "", "", "", productVendors[0].ModifiedDate, "").
					AddRow(productVendors[1].ID, "", "", productVendors[1].Name, "", "", "", "", "", "", productVendors[1].ModifiedDate, "")
		)
		defer c.db.Close()

		query := getProductVendorsByVendorQuery + " LIMIT $2 OFFSET $3"
		c.mock.ExpectQuery(query).
			WithArgs(vendorID, args.Limit, args.Offset).
			WillReturnRows(expectedResult)

		totalRows := sqlmock.NewRows([]string{"count"}).RowError(1, fmt.Errorf("row error"))
		countQuery := `SELECT COUNT(*)
			FROM product_vendor pv
			JOIN price pr ON pr.product_vendor_id = pv.id
			WHERE pr.vendor_id = $1
		`

		c.mock.ExpectQuery(countQuery).
			WithArgs(vendorID).
			WillReturnRows(totalRows)

		res, err := c.accessor.GetProductVendorsByVendor(ctx, vendorID, spec)
		c.g.Expect(err).ToNot(gomega.BeNil())
		c.g.Expect(res).To(gomega.BeNil())
	})
}

func Test_getProductByID(t *testing.T) {
	t.Parallel()

	var (
		accessor *postgresProductAccessor
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

		accessor = newPostgresProductAccessor(sqlxDB, realClock)
		mock = sqlMock

		return g, db
	}

	productFields := []string{"id", "product_category_id", "uom_id", "income_tax_id", "product_type_id", "name", "description", "modified_date", "modified_by"}

	product := &Product{
		ID:           "1",
		Name:         "Kismis",
		ModifiedDate: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(productFields).
			AddRow("1", "", "", "", "", product.Name, "", product.ModifiedDate, "")

		mock.ExpectQuery(getProductByID).
			WithArgs("1").
			WillReturnRows(rows)

		ctx := context.Background()
		res, err := accessor.getProductByID(ctx, "1")

		g.Expect(err).To(gomega.BeNil())
		g.Expect(res).To(gomega.Equal(product))
	})

	t.Run("error on row scan", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows(productFields).
			RowError(1, fmt.Errorf("error"))

		mock.ExpectQuery(getProductByID).
			WithArgs("1").
			WillReturnRows(rows)

		ctx := context.Background()
		res, err := accessor.getProductByID(ctx, "1")

		g.Expect(err).ToNot(gomega.BeNil())
		g.Expect(res).To(gomega.BeNil())
	})
}

func Test_writeProduct(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var (
			ctx     = context.Background()
			c       = setupProductAccessorTestComponent(t, WithQueryMatcher(sqlmock.QueryMatcherRegexp))
			now     = c.cmock.Now()
			product = Product{ID: "123", ModifiedDate: now}
		)

		transformedQuery, args, _ := sqlx.Named(insertProduct, product)
		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		c.mock.ExpectExec(regexp.QuoteMeta(transformedQuery)).WithArgs(
			driverArgs...,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err := c.accessor.writeProduct(ctx, product)
		c.g.Expect(err).Should(gomega.BeNil())
	})

	t.Run("error", func(t *testing.T) {
		var (
			ctx     = context.Background()
			c       = setupProductAccessorTestComponent(t)
			now     = c.cmock.Now()
			product = Product{ID: "123", ModifiedDate: now}
		)

		transformedQuery, args, _ := sqlx.Named(insertProduct, product)
		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		c.mock.ExpectExec(regexp.QuoteMeta(transformedQuery)).WithArgs(
			driverArgs...,
		).WillReturnError(errors.New("error"))

		err := c.accessor.writeProduct(ctx, product)
		c.g.Expect(err).ShouldNot(gomega.BeNil())
	})
}

func Test_writeProductCategory(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var (
			ctx      = context.Background()
			c        = setupProductAccessorTestComponent(t, WithQueryMatcher(sqlmock.QueryMatcherRegexp))
			now      = c.cmock.Now()
			category = ProductCategory{ID: "123", ModifiedDate: now}
		)

		transformedQuery, args, _ := sqlx.Named(insertProductCategory, category)
		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		c.mock.ExpectExec(regexp.QuoteMeta(transformedQuery)).WithArgs(
			driverArgs...,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err := c.accessor.writeProductCategory(ctx, category)
		c.g.Expect(err).Should(gomega.BeNil())
	})

	t.Run("error", func(t *testing.T) {
		var (
			ctx      = context.Background()
			c        = setupProductAccessorTestComponent(t)
			now      = c.cmock.Now()
			category = ProductCategory{ID: "123", ModifiedDate: now}
		)

		transformedQuery, args, _ := sqlx.Named(insertProductCategory, category)
		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		c.mock.ExpectExec(regexp.QuoteMeta(transformedQuery)).WithArgs(
			driverArgs...,
		).WillReturnError(errors.New("error"))

		err := c.accessor.writeProductCategory(ctx, category)
		c.g.Expect(err).ShouldNot(gomega.BeNil())
	})
}

func Test_writeProductType(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var (
			ctx   = context.Background()
			c     = setupProductAccessorTestComponent(t, WithQueryMatcher(sqlmock.QueryMatcherRegexp))
			now   = c.cmock.Now()
			pType = ProductType{ID: "123", ModifiedDate: now}
		)

		transformedQuery, args, _ := sqlx.Named(insertProductType, pType)
		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		c.mock.ExpectExec(regexp.QuoteMeta(transformedQuery)).WithArgs(
			driverArgs...,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err := c.accessor.writeProductType(ctx, pType)
		c.g.Expect(err).Should(gomega.BeNil())
	})

	t.Run("error", func(t *testing.T) {
		var (
			ctx   = context.Background()
			c     = setupProductAccessorTestComponent(t)
			now   = c.cmock.Now()
			pType = ProductType{ID: "123", ModifiedDate: now}
		)

		transformedQuery, args, _ := sqlx.Named(insertProductType, pType)
		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		c.mock.ExpectExec(regexp.QuoteMeta(transformedQuery)).WithArgs(
			driverArgs...,
		).WillReturnError(errors.New("error"))

		err := c.accessor.writeProductType(ctx, pType)
		c.g.Expect(err).ShouldNot(gomega.BeNil())
	})
}

func Test_writeUOM(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupProductAccessorTestComponent(t, WithQueryMatcher(sqlmock.QueryMatcherRegexp))
			now = c.cmock.Now()
			uom = UOM{ID: "123", ModifiedDate: now}
		)

		transformedQuery, args, _ := sqlx.Named(insertUOM, uom)
		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		c.mock.ExpectExec(regexp.QuoteMeta(transformedQuery)).WithArgs(
			driverArgs...,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err := c.accessor.writeUOM(ctx, uom)
		c.g.Expect(err).Should(gomega.BeNil())
	})

	t.Run("error", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupProductAccessorTestComponent(t)
			now = c.cmock.Now()
			uom = UOM{ID: "123", ModifiedDate: now}
		)

		transformedQuery, args, _ := sqlx.Named(insertUOM, uom)
		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		c.mock.ExpectExec(regexp.QuoteMeta(transformedQuery)).WithArgs(
			driverArgs...,
		).WillReturnError(errors.New("error"))

		err := c.accessor.writeUOM(ctx, uom)
		c.g.Expect(err).ShouldNot(gomega.BeNil())
	})
}

func Test_writeProductVendor(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupProductAccessorTestComponent(t, WithQueryMatcher(sqlmock.QueryMatcherRegexp))
			pv  = ProductVendor{ProductID: "123", ID: "321"}
		)

		transformedQuery, args, _ := sqlx.Named(insertProductVendor, pv)
		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		c.mock.ExpectExec(regexp.QuoteMeta(transformedQuery)).WithArgs(
			driverArgs...,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err := c.accessor.writeProductVendor(ctx, pv)
		c.g.Expect(err).Should(gomega.BeNil())
	})

	t.Run("error", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupProductAccessorTestComponent(t)
			pv  = ProductVendor{ProductID: "123", ID: "321"}
		)

		transformedQuery, args, _ := sqlx.Named(insertProductVendor, pv)
		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		c.mock.ExpectExec(regexp.QuoteMeta(transformedQuery)).WithArgs(
			driverArgs...,
		).WillReturnError(errors.New("error"))

		err := c.accessor.writeProductVendor(ctx, pv)
		c.g.Expect(err).ShouldNot(gomega.BeNil())
	})
}

func Test_writePrice(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		var (
			ctx   = context.Background()
			c     = setupProductAccessorTestComponent(t, WithQueryMatcher(sqlmock.QueryMatcherRegexp))
			price = Price{ID: "321"}
		)

		transformedQuery, args, _ := sqlx.Named(insertPrice, price)
		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		c.mock.ExpectExec(regexp.QuoteMeta(transformedQuery)).WithArgs(
			driverArgs...,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err := c.accessor.writePrice(ctx, price)
		c.g.Expect(err).Should(gomega.BeNil())
	})

	t.Run("error", func(t *testing.T) {
		var (
			ctx   = context.Background()
			c     = setupProductAccessorTestComponent(t)
			price = Price{ID: "321"}
		)

		transformedQuery, args, _ := sqlx.Named(insertPrice, price)
		driverArgs := make([]driver.Value, len(args))
		for i, arg := range args {
			driverArgs[i] = arg
		}

		c.mock.ExpectExec(regexp.QuoteMeta(transformedQuery)).WithArgs(
			driverArgs...,
		).WillReturnError(errors.New("error"))

		err := c.accessor.writePrice(ctx, price)
		c.g.Expect(err).ShouldNot(gomega.BeNil())
	})
}

func Test_Close(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		c := setupProductAccessorTestComponent(t)

		c.mock.ExpectClose()
		err := c.accessor.Close()
		c.g.Expect(err).ShouldNot(gomega.HaveOccurred())
	})

	t.Run("error", func(t *testing.T) {
		c := setupProductAccessorTestComponent(t)

		c.mock.ExpectClose().WillReturnError(errors.New("error"))
		err := c.accessor.Close()
		c.g.Expect(err).Should(gomega.HaveOccurred())
	})
}

type productAccessorTestComponent struct {
	g        *gomega.WithT
	mock     sqlmock.Sqlmock
	db       *sql.DB
	accessor *postgresProductAccessor
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

func setupProductAccessorTestComponent(t *testing.T, opts ...Option) productAccessorTestComponent {
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
		accessor: newPostgresProductAccessor(sqlxDB, clockMock),
		cmock:    clockMock,
	}
}
