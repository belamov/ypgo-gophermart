package models

import "time"

type Order struct {
	UploadedAt time.Time
	ID         int
	CreatedBy  int
	Status     OrderStatus
	Accrual    float64
}

type OrderStatus int

const (
	OrderStatusNew OrderStatus = iota + 1
	OrderStatusProcessing
	OrderStatusInvalid
	OrderStatusProcessed
)

func (s OrderStatus) String() string {
	switch s {
	case OrderStatusNew:
		return "NEW"
	case OrderStatusProcessing:
		return "PROCESSING"
	case OrderStatusInvalid:
		return "INVALID"
	case OrderStatusProcessed:
		return "PROCESSED"
	}
	return "unknown"
}
