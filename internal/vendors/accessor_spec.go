package vendors

import "kg/procurement/internal/common/database"

type AccessorGetAllPaginationData struct {
	Vendors  []Vendor                    `json:"vendors"`
	Metadata database.PaginationMetadata `json:"metadata"`
}
