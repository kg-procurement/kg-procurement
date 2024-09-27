package vendors

import (
	"context"
	"database/sql"
	"kg/procurement/internal/common/database"
	"testing"
	"time"
)

func TestVendorAccessor_GetAll(t *testing.T) {
	db, err := sql.Open("postgres", ":memory:")
	if err != nil {
		t.Fatalf("Opening in-memory database error: %w", err)
	}

	defer db.Close()

	create_table_query := `CREATE TABLE "vendor" (
  "id" int PRIMARY KEY,
  "name" varchar,
  "bp_id" int,
  "bp_name" varchar,
  "rating"  int,
  "area_group_id" int,
  "area_group_name" varchar,
  "sap_code"	  varchar,
  "modified_date" timestamp,
  "modified_by" int,
  "dt" date
);`

	_, err = db.Exec(create_table_query)
	if err != nil {
		t.Fatalf("Failed creating table vendors: %w", err)
	}

	testData := Vendor{
		Id:            1,
		Name:          "name",
		BpId:          1,
		BpName:        "bp_name",
		Rating:        1,
		AreaGroupId:   1,
		AreaGroupName: "group_name",
		SapCode:       "sap_code",
		ModifiedDate:  time.Now(),
		ModifiedBy:    1,
		Date:          time.Now(),
	}
	insertQuery := `
	INSERT INTO vendors (id, name, bp_id, bp_name, rating, area_group_id, area_group_name, sap_code, modified_date, modified_by, date)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	_, err = db.Exec(insertQuery,
		testData,
		testData.Name,
		testData.BpId,
		testData.BpName,
		testData.Rating,
		testData.AreaGroupId,
		testData.AreaGroupName,
		testData.SapCode,
		testData.ModifiedDate,
		testData.ModifiedBy,
		testData.Date,
	)

	if err != nil {
		t.Fatalf("Error while doing insert: %w", err)
	}

	accessor := newPostgresVendorAccessor()

	vendors, err := accessor.GetAll(context.Background())
}
