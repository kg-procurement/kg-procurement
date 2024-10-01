package vendors

type ServiceGetAllPaginationSpec struct {
	Limit int
	Order string
	Page  int
}

type ServiceGetAllPaginationData struct {
	Vendors      []Vendor
	TotalEntries int
	CurrentPage  int
	PreviousPage *int
	NextPage     int
}
