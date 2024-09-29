package inventory

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"testing"
	"time"

	"kg/procurement/internal/vendors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/onsi/gomega"
)

var mockVendor = vendors.Vendor{
	ID:            "v1",
	Name:          "Vendor A",
	Description:   "Test Vendor A",
	BpID:          "bp1",
	BpName:        "BP Vendor A",
	Rating:        5,
	AreaGroupID:   "ag1",
	AreaGroupName: "Area Group 1",
	SapCode:       "sap1",
	ModifiedDate:  time.Now(),
	ModifiedBy:    1,
	Date:          time.Now(),
}

func Test_newPostgresInventoryAccessor(t *testing.T) {
	_ = newPostgresInventoryAccessor(nil)
}


func TestUpdateInventory(t *testing.T) {
	t.Parallel()

	var (
		accessor *postgresInventoryAccessor
		mock     sqlmock.Sqlmock
	)

	setup := func(t *testing.T) (*gomega.GomegaWithT, *sql.DB) {
		g := gomega.NewWithT(t)
		db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			log.Fatal("error initializing mock:", err)
		}

		accessor = newPostgresInventoryAccessor(db)
		mock = sqlMock

		return g, db
	}

	t.Run("success", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		initialInventory := Item{
			ID:          "inv1",
			Name:        "Item A",
			Description: "Initial Description",
			Vendor:      mockVendor,
			Price:       100.00,
		}
		
		updatedInventory := Item{
			ID:          "inv1",
			Name:        "Item A Updated",
			Description: "Updated Description",
			Vendor:      mockVendor,
			Price:       150.00,
		}

		mock.ExpectExec(`UPDATE inventory SET name = $1, description = $2, price = $3 WHERE id = $4`).
			WithArgs(updatedInventory.Name, updatedInventory.Description, updatedInventory.Price, updatedInventory.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		ctx := context.Background()
		err := accessor.UpdateInventory(ctx, initialInventory.ID, updatedInventory)

		g.Expect(err).To(gomega.BeNil())
	})

	t.Run("error on update", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		updatedInventory := Item{
			ID:          "inv1",
			Name:        "Item A Updated",
			Description: "Updated Description",
			Vendor:      mockVendor,
			Price:       150.00,
		}

		mock.ExpectExec(`UPDATE inventory SET name = $1, description = $2, price = $3 WHERE id = $4`).
			WithArgs(updatedInventory.Name, updatedInventory.Description, updatedInventory.Price, updatedInventory.ID).
			WillReturnError(errors.New("some database error"))

		ctx := context.Background()
		err := accessor.UpdateInventory(ctx, updatedInventory.ID, updatedInventory)

		g.Expect(err).ToNot(gomega.BeNil())
	})

	t.Run("inventory not found", func(t *testing.T) {
		g, db := setup(t)
		defer db.Close()

		updatedInventory := Item{
			ID:          "nonexistent_inv",
			Name:        "Non-existent Item",
			Description: "This item does not exist",
			Vendor:      mockVendor,
			Price:       0.00,
		}

		mock.ExpectExec(`UPDATE inventory SET name = $1, description = $2, price = $3 WHERE id = $4`).
			WithArgs(updatedInventory.Name, updatedInventory.Description, updatedInventory.Price, updatedInventory.ID).
			WillReturnError(errors.New("inventory item not found"))

		ctx := context.Background()
		err := accessor.UpdateInventory(ctx, updatedInventory.ID, updatedInventory)

		g.Expect(err).ToNot(gomega.BeNil())
	})
}
