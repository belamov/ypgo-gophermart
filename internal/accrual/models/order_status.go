package models

const (
	OrderStatusNew OrderStatus = iota + 1
	OrderStatusProcessing
	OrderStatusInvalid
	OrderStatusProcessed
	OrderStatusError
)

type OrderStatus int

func (s OrderStatus) String() string {
	switch s {
	case OrderStatusNew:
		return "REGISTERED"
	case OrderStatusProcessing:
		return "PROCESSING"
	case OrderStatusInvalid:
		return "INVALID"
	case OrderStatusProcessed:
		return "PROCESSED"
	case OrderStatusError:
		return "PROCESSING"
	}
	return "unknown"
}
