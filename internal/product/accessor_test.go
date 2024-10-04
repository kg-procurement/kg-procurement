package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/onsi/gomega"
)

func Test_newPostgresProductAccessor(t *testing.T) {
	_ = newPostgresProductAccessor(nil)
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

	query := `UPDATE product
        SET 
            product_category_id = $2,
            uom_id = $3,
            income_tax_id = $4,
            product_type_id = $5,
            name = $6,
            description = $7,
            modified_date = $8,
            modified_by = $9
        WHERE 
            id = $1
        RETURNING 
            id,
            product_category_id,
            uom_id,
            income_tax_id,
            product_type_id,
            name,
            description,
            modified_date,
            modified_by
    `

	fixedTime := time.Now()
	updatedFixedTime := time.Now().Add(1 * time.Hour)

	t.Run("success", func(t *testing.T) {
		var (
            ctx            = context.Background()
            c              = setupProductAccessorTestComponent(t)
            expectedResult = sqlmock.NewRows(productFields).
                AddRow(
                    "inv1",
                    "category_id_updated",
                    "uom_id_updated",
                    "income_tax_id_updated",
                    "product_type_id_updated",
                    "Product A Updated",
                    "Updated Description",
                    updatedFixedTime,
                    "modified_by_updated",
                )
        )
        defer c.db.Close()

		updatedProduct := &Product{
            ID:                "inv1",
            ProductCategoryID: "category_id_updated",
            UOMID:             "uom_id_updated",
            IncomeTaxID:       "income_tax_id_updated",
            ProductTypeID:     "product_type_id_updated",
            Name:              "Product A Updated",
            Description:       "Updated Description",
            ModifiedDate:      updatedFixedTime,
            ModifiedBy:        "modified_by_updated",
        }

		c.mock.ExpectQuery(query).
            WithArgs(
                updatedProduct.ID,
                updatedProduct.ProductCategoryID,
                updatedProduct.UOMID,
                updatedProduct.IncomeTaxID,
                updatedProduct.ProductTypeID,
                updatedProduct.Name,
                updatedProduct.Description,
                updatedProduct.ModifiedDate,
                updatedProduct.ModifiedBy,
            ).WillReturnRows(expectedResult)

        res, err := c.accessor.UpdateProduct(ctx, *updatedProduct)

        c.g.Expect(err).To(gomega.BeNil())
        c.g.Expect(res).To(gomega.Equal(updatedProduct))
	})

	t.Run("error on row scan", func(t *testing.T) {
        var (
            ctx            = context.Background()
            c              = setupProductAccessorTestComponent(t)
            expectedResult = sqlmock.NewRows(productFields).
                AddRow(
                    "inv1",
                    "category_id_updated",
                    "uom_id_updated",
                    "income_tax_id_updated",
                    "product_type_id_updated",
                    "Product A Updated",
                    "Updated Description",
                    "not a time",
                    "modified_by_updated",
                )
        )
        defer c.db.Close()

        c.mock.ExpectQuery(query).
            WithArgs(
                "inv1",
                "category_id_updated",
                "uom_id_updated",
                "income_tax_id_updated",
                "product_type_id_updated",
                "Product A Updated",
                "Updated Description",
                "not a time",
                "modified_by_updated",
            ).WillReturnRows(expectedResult)

        updatedProduct := &Product{
            ID:                "inv1",
            ProductCategoryID: "category_id_updated",
            UOMID:             "uom_id_updated",
            IncomeTaxID:       "income_tax_id_updated",
            ProductTypeID:     "product_type_id_updated",
            Name:              "Product A Updated",
            Description:       "Updated Description",
            ModifiedDate:      fixedTime,
            ModifiedBy:        "modified_by_updated",
        }

        res, err := c.accessor.UpdateProduct(ctx, *updatedProduct)

        c.g.Expect(err).ToNot(gomega.BeNil())
        c.g.Expect(res).To(gomega.BeNil())
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
            modified_date = $25,
            modified_by = $26
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

	// Example time values for your test
	fixedTime := time.Now()
	updatedFixedTime := time.Now().Add(1 * time.Hour)

	t.Run("success", func(t *testing.T) {
		var (
			ctx            = context.Background()
			c              = setupProductAccessorTestComponent(t)
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
					updatedFixedTime,
					"modified_by_updated",
				)
		)
		defer c.db.Close()

		updatedPrice := &Price{
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
			ModifiedDate:    updatedFixedTime,
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
				updatedPrice.ModifiedDate,
				updatedPrice.ModifiedBy,
			).WillReturnRows(expectedResult)

		res, err := c.accessor.UpdatePrice(ctx, *updatedPrice)

		c.g.Expect(err).To(gomega.BeNil())
		c.g.Expect(res).To(gomega.Equal(updatedPrice))
	})

	t.Run("error on row scan", func(t *testing.T) {
		var (
			ctx            = context.Background()
			c              = setupProductAccessorTestComponent(t)
			expectedResult = sqlmock.NewRows(priceFields).
					AddRow(
					"ID",
					"org_id_updated",
					"vendor_id_updated",
					"product_vendor_id_updated",
					10,
					100,
					"quantity_uom_id_updated",
					"not int", // Invalid quantity
					"not int", // Invalid quantity
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
					"modified_by_updated",
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
				"modified_by_updated",
			).WillReturnRows(expectedResult)

		updatedPrice := &Price{
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
			ModifiedDate:    updatedFixedTime,
			ModifiedBy:      "modified_by_updated",
		}

		res, err := c.accessor.UpdatePrice(ctx, *updatedPrice)

		c.g.Expect(err).ToNot(gomega.BeNil())
		c.g.Expect(res).To(gomega.BeNil())
	})
}

func Test_GetProductsByVendor(t *testing.T) {
	t.Parallel()

	var (
		vendorID       = "1234"
		productColumns = []string{"id", "product_category_id", "uom_id", "income_tax_id", "product_type_id", "name", "description", "modified_date", "modified_by"}
		products       = []Product{
			{
				ID:           "1111",
				Name:         "Mixer",
				ModifiedDate: time.Now(),
			},
			{
				ID:           "2222",
				Name:         "Rice Cooker",
				ModifiedDate: time.Now(),
			},
		}
	)

	t.Run("success", func(t *testing.T) {
		var (
			ctx            = context.Background()
			c              = setupProductAccessorTestComponent(t)
			expectedResult = sqlmock.NewRows(productColumns).
					AddRow(products[0].ID, "", "", "", "", products[0].Name, "", products[0].ModifiedDate, "").
					AddRow(products[1].ID, "", "", "", "", products[1].Name, "", products[1].ModifiedDate, "")
		)
		defer c.db.Close()

		c.mock.ExpectQuery(getProductsByVendorQuery).
			WithArgs(vendorID).
			WillReturnRows(expectedResult)

		res, err := c.accessor.GetProductsByVendor(ctx, vendorID)

		c.g.Expect(err).To(gomega.BeNil())
		c.g.Expect(res).To(gomega.BeComparableTo(products))
	})

	t.Run("error on query execution", func(t *testing.T) {
		var (
			ctx = context.Background()
			c   = setupProductAccessorTestComponent(t)
		)
		defer c.db.Close()

		c.mock.ExpectQuery(getProductsByVendorQuery).
			WithArgs(vendorID).
			WillReturnError(errors.New("error"))

		res, err := c.accessor.GetProductsByVendor(ctx, vendorID)

		c.g.Expect(err).ShouldNot(gomega.BeNil())
		c.g.Expect(res).To(gomega.BeNil())
	})

	t.Run("error on scanning data row", func(t *testing.T) {
		var (
			ctx            = context.Background()
			c              = setupProductAccessorTestComponent(t)
			expectedResult = sqlmock.NewRows(productColumns).
					AddRow(products[0].ID, "", "", "", "", products[0].Name, "", products[0].ModifiedDate, "").
					AddRow(products[1].ID, "", "", "", "", products[1].Name, "", products[1].ModifiedDate, "").
					AddRow(nil, nil, nil, nil, nil, nil, nil, nil, nil)
		)
		defer c.db.Close()

		c.mock.ExpectQuery(getProductsByVendorQuery).
			WithArgs(vendorID).
			WillReturnRows(expectedResult)

		res, err := c.accessor.GetProductsByVendor(ctx, vendorID)

		c.g.Expect(err).ShouldNot(gomega.BeNil())
		c.g.Expect(res).To(gomega.BeNil())
	})

	t.Run("error while iterating rows", func(t *testing.T) {
		var (
			ctx            = context.Background()
			c              = setupProductAccessorTestComponent(t)
			expectedResult = sqlmock.NewRows(productColumns).
					AddRow(products[0].ID, "", "", "", "", products[0].Name, "", products[0].ModifiedDate, "").
					AddRow(products[1].ID, "", "", "", "", products[1].Name, "", products[1].ModifiedDate, "").
					RowError(0, fmt.Errorf("some error"))
		)
		defer c.db.Close()

		c.mock.ExpectQuery(getProductsByVendorQuery).
			WithArgs(vendorID).
			WillReturnRows(expectedResult)

		res, err := c.accessor.GetProductsByVendor(ctx, vendorID)

		c.g.Expect(err).ToNot(gomega.BeNil())
		c.g.Expect(res).To(gomega.BeNil())
	})
}

type productAccessorTestComponent struct {
	g        *gomega.WithT
	mock     sqlmock.Sqlmock
	db       *sql.DB
	accessor *postgresProductAccessor
}

func setupProductAccessorTestComponent(t *testing.T) productAccessorTestComponent {
	g := gomega.NewWithT(t)
	db, sqlMock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	return productAccessorTestComponent{
		g:        g,
		mock:     sqlMock,
		db:       db,
		accessor: newPostgresProductAccessor(db),
	}
}
