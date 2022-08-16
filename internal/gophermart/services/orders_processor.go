package services

import (
	"context"
	"errors"
	"sync"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/storage"
	"github.com/cenkalti/backoff/v4"
	"github.com/rs/zerolog/log"
)

type OrderProcessor struct {
	OrdersStorage     storage.OrdersStorage
	AccrualService    AccrualInfoProvider
	BalanceManager    BalanceProcessorInterface
	ordersToProcessCh chan models.Order
	stopCh            chan struct{}
	wg                sync.WaitGroup
}

func NewOrderProcessor(ordersStorage storage.OrdersStorage, accrualService AccrualInfoProvider, balanceManager BalanceProcessorInterface) *OrderProcessor {
	return &OrderProcessor{
		OrdersStorage:     ordersStorage,
		AccrualService:    accrualService,
		BalanceManager:    balanceManager,
		ordersToProcessCh: make(chan models.Order),
		stopCh:            make(chan struct{}),
		wg:                sync.WaitGroup{},
	}
}

func (o *OrderProcessor) RegisterOrderForProcessing(order models.Order) {
	for {
		select {
		// if stopCh is closed, we will not be sending orders anymore
		// stopCh closed only in func that reads from ordersToProcessCh
		case <-o.stopCh:
			log.Debug().Int("order_id", order.ID).Msg("ignoring order processing, received stop signal")
			return
		case o.ordersToProcessCh <- order:
			log.Debug().Int("order_id", order.ID).Msg("registering order for processing")
			return
		}
	}
}

func (o *OrderProcessor) StartProcessing(ctx context.Context) {
	newOrders, err := o.OrdersStorage.GetOrdersForProcessing()
	if err != nil {
		log.Error().
			Err(err).
			Msg("cant fetch orders for processing")
	}

	for _, newOrder := range newOrders {
		log.Debug().Int("order_io", newOrder.ID).Msg("start processing order")
		go o.ProcessOrder(ctx, newOrder)
	}

	for {
		select {
		case <-ctx.Done():
			// when context is canceled, we will signal all senders that they should stop
			// sending orders to channel, because we won't process them anymore
			close(o.stopCh)
			log.Debug().Msg("stop receiving orders to process")
			return
		case orderToProcess := <-o.ordersToProcessCh:
			log.Debug().Int("order_id", orderToProcess.ID).Msg("received order to process")
			go o.ProcessOrder(ctx, orderToProcess)

		}
	}
}

func (o *OrderProcessor) ProcessOrder(ctx context.Context, order models.Order) {
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
	backOff := backoff.WithContext(exponentialBackOff, ctx)

	addAccrualForOrder := func() error {
		accrual, err := o.AccrualService.GetAccrualForOrder(ctx, order.ID)

		// order is not yet proceeded, we will try to fetch it later
		if errors.Is(err, ErrOrderIsNotYetProceeded) {
			return err
		}

		// order is proceeded, but is invalid. we will try to fetch it later, when app is restarted
		if errors.Is(err, ErrInvalidOrderForAccrual) {
			err := o.OrdersStorage.ChangeStatus(order, models.OrderStatusNew)
			if err != nil {
				log.Error().
					Err(err).
					Int("order_id", order.ID).
					Int("new_order_status", int(models.OrderStatusNew)).
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
		log.Debug().
			Int("order_id", order.ID).
			Float64("accrual", accrual).
			Msg("fetched info about order accrual. received accrual")
		err = o.BalanceManager.AddAccrual(order, accrual)
		if err != nil {
			return backoff.Permanent(err)
		}

		return nil
	}

	err = backoff.Retry(addAccrualForOrder, backOff)
	if err != nil {
		if errors.Is(err, ctx.Err()) {
			log.Info().
				Int("order_id", order.ID).
				Msg("canceling order processing gracefully. returning order to NEW status")
		} else {
			log.Error().
				Err(err).
				Int("order_id", order.ID).
				Msg("received unexpected error while processing order. returning order to NEW status")
		}

		err := o.OrdersStorage.ChangeStatus(order, models.OrderStatusNew)
		if err != nil {
			log.Error().
				Err(err).
				Int("order_id", order.ID).
				Int("new_order_status", int(models.OrderStatusNew)).
				Msg("unexpected error while processing order. cant change order status")
			return
		}
	}
}
