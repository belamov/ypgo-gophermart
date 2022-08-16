package models

import "time"

type Order struct {
	ID        int
	CreatedAt time.Time
	Status    OrderStatus
	Accrual   float64
	Items     []OrderItem
}
