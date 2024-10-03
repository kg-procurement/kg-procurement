package product

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

func Test_newPostgresProductAccessor(t *testing.T) {
	_ = newPostgresProductAccessor(nil)
}

func TestUpdateProduct(t *testing.T) {
	t.Parallel()

	var (
		accessor *postgresProductAccessor
		mock     sqlmock.Sqlmock
	)

	setup := func(t *testing.T) (*gomega.GomegaWithT, *sql.DB) {
		g := gomega.NewWithT(t)
		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			log.Fatal("error initializing mock:", err)
		}

		accessor = newPostgresProductAccessor(db)
		mock = sqlMock

		return g, db
	}

	t.Run("When updating product successfully, should return no error", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		initialProduct := Product{
			ID:          "inv1",
			Name:        "Product A",
			Description: "Initial Description",
		}

		updatedProduct := Product{
			ID:          "inv1",
			Name:        "Product A Updated",
			Description: "Updated Description",
		}

		mock.ExpectExec(`UPDATE product SET name = $1 description = $2 WHERE id = $3`).
			WithArgs(updatedProduct.Name, updatedProduct.Description, updatedProduct.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		ctx := context.Background()
		err := accessor.UpdateProduct(ctx, initialProduct.ID, updatedProduct)

		g.Expect(err).To(gomega.BeNil())
	})

	t.Run("When updating product fails due to database error, should return error", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		updatedProduct := Product{
			ID:          "inv1",
			Name:        "Product A Updated",
			Description: "Updated Description",
		}

		mock.ExpectExec(`UPDATE product SET name = $1 description = $2 WHERE id = $3`).
			WithArgs(updatedProduct.Name, updatedProduct.Description, updatedProduct.ID).
			WillReturnError(errors.New("some database error"))

		ctx := context.Background()
		err := accessor.UpdateProduct(ctx, updatedProduct.ID, updatedProduct)

		g.Expect(err).ToNot(gomega.BeNil())
	})

	t.Run("When updating product and item not found, should return error", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		updatedProduct := Product{
			ID:          "nonexistent_inv",
			Name:        "Non-existent Product",
			Description: "This item does not exist",
		}

		mock.ExpectExec(`UPDATE product SET name = $1 description = $2 WHERE id = $3`).
			WithArgs(updatedProduct.Name, updatedProduct.Description, updatedProduct.ID).
			WillReturnError(errors.New("product item not found"))

		ctx := context.Background()
		err := accessor.UpdateProduct(ctx, updatedProduct.ID, updatedProduct)

		g.Expect(err).ToNot(gomega.BeNil())
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
