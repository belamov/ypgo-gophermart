package models

type OrderStatus int

const (
	OrderStatusNew OrderStatus = iota + 1
	OrderStatusProcessing
	OrderStatusInvalid
	OrderStatusProcessed
	OrderStatusError
)
