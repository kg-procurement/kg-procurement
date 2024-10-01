package vendors

import "kg/procurement/internal/common/database"

type ServiceGetAllPaginationSpec struct {
	Limit int
	Order string
	Page  int
}

type ServiceGetAllPaginationData struct {
	Vendors  []Vendor
	Metadata database.PaginationMetadata
}
