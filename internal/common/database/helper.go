package database

import (
	"strings"
)

type PaginationSpec struct {
	Limit   int    `json:"limit"`
	Page    int    `json:"page"`
	Order   string `json:"order"`
	OrderBy string `json:"order_by"`
}

type PaginationArgs struct {
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
	Order   string `json:"order"`
	OrderBy string `json:"order_by"`
}

type PaginationMetadata struct {
	TotalPage    int `json:"total_page"`
	CurrentPage  int `json:"current_page"`
	TotalEntries int `json:"total_entries"`
}

func validateOrderString(order string) string {
	uppercaseOrder := strings.ToUpper(order)
	switch uppercaseOrder {
	case "ASC":
		return uppercaseOrder
	case "DESC":
		return uppercaseOrder
	default:
		return "ASC"
	}
}

func BuildPaginationArgs(spec PaginationSpec) PaginationArgs {
	limit := 10
	if spec.Limit > 0 {
		limit = spec.Limit
	}

	offset := spec.Limit * (spec.Page - 1)

	return PaginationArgs{
		Limit:   limit,
		Offset:  offset,
		Order:   validateOrderString(spec.Order),
		OrderBy: spec.OrderBy,
	}
}

func GeneratePaginationMetadata(spec PaginationSpec, totalEntries int) PaginationMetadata {
	totalPage := (totalEntries / spec.Limit)
	if totalEntries%spec.Limit != 0 {
		totalPage += 1
	}

	return PaginationMetadata{
		TotalPage:    totalPage,
		CurrentPage:  spec.Page,
		TotalEntries: totalEntries,
	}

}
