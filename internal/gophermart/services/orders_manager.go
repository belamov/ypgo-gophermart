package services

import (
	"errors"

	"github.com/belamov/ypgo-gophermart/internal"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/storage"
)

type OrdersManagerInterface interface {
	AddOrder(orderID int, userID int) error
	ValidateOrderID(s int) error
	GetUsersOrders(userID int) ([]models.Order, error)
	ProcessOrder(order models.Order)
}

type OrdersManager struct {
	OrdersStorage    storage.OrdersStorage
	BalanceProcessor BalanceProcessorInterface
	OrdersProcessor  *OrderProcessor
}

func (o *OrdersManager) GetUsersOrders(userID int) ([]models.Order, error) {
	return o.OrdersStorage.GetUsersOrders(userID)
}

func NewOrdersManager(ordersStorage storage.OrdersStorage, balanceProcessor BalanceProcessorInterface, ordersProcessor *OrderProcessor) *OrdersManager {
	return &OrdersManager{
		ordersStorage,
		balanceProcessor,
		ordersProcessor,
	}
}

func (o *OrdersManager) AddOrder(orderID int, userID int) error {
	existingOrder, err := o.OrdersStorage.FindByID(orderID)
	if err != nil {
		return err
	}

	if existingOrder.ID != 0 {
		return NewOrderAlreadyAddedError(existingOrder)
	}

	createdOrder, err := o.OrdersStorage.CreateNew(orderID, userID)
	if err != nil {
		return err
	}

	go o.ProcessOrder(createdOrder)

	return nil
}

func (o *OrdersManager) ValidateOrderID(orderID int) error {
	if orderID <= 0 {
		return errors.New("order ID should be greater than zero")
	}

	if !internal.ValidLuhn(orderID) {
		return errors.New("order ID is not validated by Luhn algorithm")
	}

	return nil
}

// ProcessOrder sends order to queue for processing. Blocks until processing is begun
func (o *OrdersManager) ProcessOrder(order models.Order) {
	o.OrdersProcessor.RegisterOrderForProcessing(order)
}
