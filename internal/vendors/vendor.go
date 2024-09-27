package vendors

import "time"

// Vendor defines the metadata related to a vendor
// i.e. name, etc
type Vendor struct {
	Id            int
	Name          string
	BpId          int
	BpName        string
	Rating        int
	AreaGroupId   int
	AreaGroupName string
	SapCode       string
	ModifiedDate  time.Time
	ModifiedBy    int
	Date          time.Time
}
