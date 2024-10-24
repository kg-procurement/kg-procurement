package vendors

import (
	"kg/procurement/internal/common/database"
	"time"
)

// Vendor defines the metadata related to a vendor
// i.e. name, etc
type Vendor struct {
	ID            string    `db:"id" json:"id"`
	Email         string    `db:"email" json:"email"`
	Name          string    `db:"name" json:"name"`
	Description   string    `db:"description" json:"description"`
	BpID          string    `db:"bp_id" json:"bp_id"`
	BpName        string    `db:"bp_name" json:"bp_name"`
	Rating        int       `db:"rating" json:"rating"`
	AreaGroupID   string    `db:"area_group_id" json:"area_group_id"`
	AreaGroupName string    `db:"area_group_name" json:"area_group_name"`
	SapCode       string    `db:"sap_code" json:"sap_code"`
	ModifiedDate  time.Time `db:"modified_date" json:"modified_date"`
	ModifiedBy    string    `db:"modified_by" json:"modified_by"`
	Date          time.Time `db:"dt" json:"dt"`
}

type PutVendorSpec struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	BpID          string `json:"bp_id"`
	BpName        string `json:"bp_name"`
	Rating        int    `json:"rating"`
	AreaGroupID   string `json:"area_group_id"`
	AreaGroupName string `json:"area_group_name"`
	SapCode       string `json:"sap_code"`
}

type GetAllVendorSpec struct {
	Location string `json:"location"`
	Product  string `json:"product"`
	database.PaginationSpec
}
