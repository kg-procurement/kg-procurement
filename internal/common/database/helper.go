package database

import (
	"strings"
)

type PaginationSpec struct {
	Limit   int
	Page    int
	Order   string
	OrderBy string
}

type PaginationArgs struct {
	Limit   int
	Offset  int
	Order   string
	OrderBy string
}

type PaginationMetadata struct {
	TotalPage    int
	CurrentPage  int
	TotalEntries int
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
