package services

type OrdersProcessorInterface interface {
	AddOrder(orderID int, userID int) error
	ValidateOrderId(s int) error
}

type OrdersProcessor struct{}

func (o *OrdersProcessor) AddOrder(orderID int, userID int) error {
	// TODO implement me
	panic("implement me")
}

func (o *OrdersProcessor) ValidateOrderId(s int) error {
	// TODO implement me
	panic("implement me")
}

func NewOrdersProcessor() *OrdersProcessor {
	return &OrdersProcessor{}
}
