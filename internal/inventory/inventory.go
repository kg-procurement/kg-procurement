package inventory

import (
	"kg/procurement/internal/vendors"
)

type Item struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Vendor      vendors.Vendor `json:"vendor"`
	Price       float64        `json:"price"`
}
