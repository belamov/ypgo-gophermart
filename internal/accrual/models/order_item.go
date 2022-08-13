package models

type OrderItem struct {
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
