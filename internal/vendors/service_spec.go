package vendors

import "kg/procurement/internal/common/database"

type ServiceGetAllPaginationData struct {
	Vendors  []Vendor
	Metadata database.PaginationMetadata
}
