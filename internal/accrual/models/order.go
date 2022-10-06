package models

import "time"

type Order struct {
	CreatedAt time.Time
	Items     []OrderItem
	ID        int
	Status    OrderStatus
	Accrual   float64
}
