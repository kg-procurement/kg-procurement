package router

import (
	"kg/procurement/internal/common/database"
	"net/http"
	"strconv"
)

func GetPaginationSpec(r *http.Request) database.PaginationSpec {
	queryParam := r.URL.Query()

	pageString := queryParam.Get("page")
	page, err := strconv.Atoi(pageString)
	if err != nil {
		page = 1
	}

	order := queryParam.Get("order")

	spec := database.PaginationSpec{
		Limit: 10,
		Order: order,
		Page:  page,
	}

	return spec
}
