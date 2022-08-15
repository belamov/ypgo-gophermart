package services

import (
	"context"
	"errors"

	"github.com/belamov/ypgo-gophermart/internal"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/storage"
	"github.com/cenkalti/backoff/v4"
	"github.com/rs/zerolog/log"
)

type OrdersProcessorInterface interface {
	AddOrder(orderID int, userID int) error
	ValidateOrderID(s int) error
	GetUsersOrders(userID int) ([]models.Order, error)
}

type OrdersProcessor struct {
	OrdersStorage    storage.OrdersStorage
	BalanceProcessor BalanceProcessorInterface
	AccrualService   AccrualInfoProvider
}

func (o *OrdersProcessor) GetUsersOrders(userID int) ([]models.Order, error) {
	return o.OrdersStorage.GetUsersOrders(userID)
}

func NewOrdersProcessor(ordersStorage storage.OrdersStorage, balanceProcessor BalanceProcessorInterface, accrualService AccrualInfoProvider) *OrdersProcessor {
	return &OrdersProcessor{
		ordersStorage,
		balanceProcessor,
		accrualService,
	}
}

func (o *OrdersProcessor) AddOrder(orderID int, userID int) error {
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

func (o *OrdersProcessor) ValidateOrderID(orderID int) error {
	if orderID <= 0 {
		return errors.New("order ID should be greater than zero")
	}

	if !internal.ValidLuhn(orderID) {
		return errors.New("order ID is not validated by Luhn algorithm")
	}

	return nil
}

func (o *OrdersProcessor) ProcessOrder(order models.Order) {
	err := o.OrdersStorage.ChangeStatus(order, models.OrderStatusProcessing)
	if err != nil {
		log.Error().
			Err(err).
			Int("order_id", order.ID).
			Int("new_order_status", int(models.OrderStatusProcessing)).
			Msg("unexpected error while processing order. cant change order status")
		return
	}

	exponentialBackOff := backoff.NewExponentialBackOff()
	exponentialBackOff.MaxElapsedTime = 0
	backOff := backoff.WithContext(exponentialBackOff, context.Background()) // todo: external context
	orderProcessOperation := func() error {
		accrual, err := o.AccrualService.GetAccrualForOrder(context.Background(), order.ID) // todo: external context

		// order is not yet proceeded, we will try to fetch it later
		if errors.Is(err, ErrOrderIsNotYetProceeded) {
			return err
		}

		// order is proceeded, but no accrual will be added in future
		if errors.Is(err, ErrInvalidOrderForAccrual) {
			err := o.OrdersStorage.ChangeStatus(order, models.OrderStatusInvalid)
			if err != nil {
				log.Error().
					Err(err).
					Int("order_id", order.ID).
					Int("new_order_status", int(models.OrderStatusInvalid)).
					Msg("unexpected error while processing order. cant change order status")
				return backoff.Permanent(err)
			}
			return nil
		}

		// unexpected error
		if err != nil {
			log.Error().
				Err(err).
				Int("order_id", order.ID).
				Msg("received unexpected error from accrual service")
			return backoff.Permanent(err)
		}

		// order is proceeded, accrual is available
		err = o.BalanceProcessor.AddAccrual(order, accrual)
		if err != nil {
			return backoff.Permanent(err)
		}

		return nil
	}

	err = backoff.Retry(orderProcessOperation, backOff)
	if err != nil {
		log.Error().
			Err(err).
			Int("order_id", order.ID).
			Msg("received unexpected error while trying to fetch accrual")
	}
}
