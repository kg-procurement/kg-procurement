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
	orderBy := queryParam.Get("order_by")

	limitStr := queryParam.Get("limit")
	limit := 10 // Default value
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			limit = 10
		}
	}

	spec := database.PaginationSpec{
		Limit:   limit,
		Order:   order,
		Page:    page,
		OrderBy: orderBy,
	}

	return spec
}
