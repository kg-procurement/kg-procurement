package vendors

import "kg/procurement/internal/common/database"

type AccessorGetAllPaginationData struct {
	Vendors  []Vendor
	Metadata database.PaginationMetadata
}
