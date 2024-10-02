package router

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
