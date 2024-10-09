package vendors

import (
	"kg/procurement/internal/common/database"
	"time"
)

// Vendor defines the metadata related to a vendor
// i.e. name, etc
type Vendor struct {
	ID            string    `db:"id"`
	Name          string    `db:"name"`
	Description   string    `db:"description"`
	BpID          string    `db:"bp_id"`
	BpName        string    `db:"bp_name"`
	Rating        int       `db:"rating"`
	AreaGroupID   string    `db:"area_group_id"`
	AreaGroupName string    `db:"area_group_name"`
	SapCode       string    `db:"sap_code"`
	ModifiedDate  time.Time `db:"modified_date"`
	ModifiedBy    string    `db:"modified_by"`
	Date          time.Time `db:"dt"`
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
