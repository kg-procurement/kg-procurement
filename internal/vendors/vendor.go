package vendors

import (
	"kg/procurement/internal/common/database"
	"time"
)

// Vendor defines the metadata related to a vendor
// i.e. name, etc
type Vendor struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	BpID          string    `json:"bp_id"`
	BpName        string    `json:"bp_name"`
	Rating        int       `json:"rating"`
	AreaGroupID   string    `json:"area_group_id"`
	AreaGroupName string    `json:"area_group_name"`
	SapCode       string    `json:"sap_code"`
	ModifiedDate  time.Time `json:"modified_date"`
	ModifiedBy    string    `json:"modified_by"`
	Date          time.Time `json:"dt"`
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
