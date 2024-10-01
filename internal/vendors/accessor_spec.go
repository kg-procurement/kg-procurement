package vendors

type AccessorGetAllPaginationSpec struct {
	Limit  int
	Offset int
	Order  string
}

type AccessorGetAllPaginationData struct {
	Vendors      []Vendor
	TotalEntries int
}
