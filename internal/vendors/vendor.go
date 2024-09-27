package vendors

import "time"

// Vendor defines the metadata related to a vendor
// i.e. name, etc
type Vendor struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	BpId          int       `json:"bp_id"`
	BpName        string    `json:"bp_name"`
	Rating        int       `json:"rating"`
	AreaGroupId   int       `json:"area_group_id"`
	AreaGroupName string    `json:"area_group_name"`
	SapCode       string    `json:"sap_code"`
	ModifiedDate  time.Time `json:"modified_date"`
	ModifiedBy    int       `json:"modified_by"`
	Date          time.Time `json:"dt"`
}
