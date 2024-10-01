package database

import (
	"strings"
)

type GetAllPaginationSpec struct {
	Limit int
	Page  int
	Order string
}

type GetAllPaginationArgs struct {
	Limit  int
	Offset int
	Order  string
}

type PaginationMetadata struct {
	TotalPage   int
	CurrentPage int
}

func ValidateOrderString(order string) string {
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

func GeneratePaginationArgs(spec GetAllPaginationSpec) GetAllPaginationArgs {
	limit := 10
	if spec.Limit > 0 {
		limit = spec.Limit
	}

	offset := spec.Limit * (spec.Page - 1)

	paginationArgs := GetAllPaginationArgs{
		Limit:  limit,
		Offset: offset,
		Order:  spec.Order,
	}

	return paginationArgs
}

func GeneratePaginationMetadata(spec GetAllPaginationSpec, totalEntries int) PaginationMetadata {
	totalPage := (totalEntries / spec.Limit)
	if totalEntries%spec.Limit != 0 {
		totalPage += 1
	}

	metadata := PaginationMetadata{
		TotalPage:   totalPage,
		CurrentPage: spec.Page,
	}

	return metadata
}
