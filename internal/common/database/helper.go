package database

import "strings"

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
