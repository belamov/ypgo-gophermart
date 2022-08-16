package services

import (
	"context"

	"github.com/belamov/ypgo-gophermart/internal/accrual/models"
	"github.com/belamov/ypgo-gophermart/internal/accrual/storage"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type AccrualProcessor struct {
	ordersStorage     storage.OrdersStorage
	rewardsStorage    storage.RewardsStorage
	ordersToProcessCh chan models.Order
	stopCh            chan struct{}
}

func NewAccrualProcessor(ordersStorage storage.OrdersStorage, rewardsStorage storage.RewardsStorage) *AccrualProcessor {
	return &AccrualProcessor{
		ordersStorage:     ordersStorage,
		rewardsStorage:    rewardsStorage,
		ordersToProcessCh: make(chan models.Order),
		stopCh:            make(chan struct{}),
	}
}

func (p *AccrualProcessor) RegisterOrderForProcessing(order models.Order) {
	for {
		select {
		// if stopCh is closed, we will not be sending orders anymore
		// stopCh closed only in func that reads from ordersToProcessCh
		case <-p.stopCh:
			log.Debug().Int("order_id", order.ID).Msg("ignoring order processing, received stop signal")
			return
		case p.ordersToProcessCh <- order:
			log.Debug().Int("order_id", order.ID).Msg("registering order for processing")
			return
		}
	}
}

func (p *AccrualProcessor) StartProcessing(ctx context.Context) {
	newOrders, err := p.ordersStorage.GetOrdersForProcessing()
	if err != nil {
		log.Error().
			Err(err).
			Msg("cant fetch orders for processing")
	}

	for _, newOrder := range newOrders {
		log.Debug().Int("order_id", newOrder.ID).Msg("start processing order")
		go p.AddAccrualForOrder(ctx, newOrder)
	}

	for {
		select {
		case <-ctx.Done():
			// when context is canceled, we will signal all senders that they should stop
			// sending orders to channel, because we won't process them anymore
			close(p.stopCh)
			log.Debug().Msg("stop receiving orders to process")
			return
		case orderToProcess := <-p.ordersToProcessCh:
			log.Debug().Int("order_id", orderToProcess.ID).Msg("received order to process")
			go p.AddAccrualForOrder(ctx, orderToProcess)

		}
	}
}

func (p *AccrualProcessor) AddAccrualForOrder(ctx context.Context, order models.Order) {
	err := p.ordersStorage.ChangeStatus(order.ID, models.OrderStatusProcessing)
	if err != nil {
		log.Error().Err(err).Int("order_id", order.ID).Msg("unexpected error while changing status of order")
		p.markOrderAsFailed(order)
		return
	}

	g, ctx := errgroup.WithContext(ctx)

	results := make([]float64, len(order.Items))

	log.Debug().Int("order_id", order.ID).Msg("calculating accrual for order")

	for i, orderItem := range order.Items {
		i, orderItem := i, orderItem
		g.Go(func() error {
			result, err := p.calculateAccrualForOrderItem(ctx, orderItem)
			if err == nil {
				results[i] = result
			}
			return err
		})
	}

	if err := g.Wait(); err != nil {
		log.Error().
			Err(err).
			Int("order_id", order.ID).
			Msg("unexpected error while adding accrual to order. some of order items not calculated properly")
		p.markOrderAsFailed(order)
		return
	}

	totalAccrual := 0.0
	for _, accrualForItem := range results {
		totalAccrual += accrualForItem
	}

	log.Debug().Int("order_id", order.ID).Float64("accrual", totalAccrual).Msg("calculated accrual for order")

	err = p.ordersStorage.AddAccrual(order.ID, totalAccrual)
	if err != nil {
		log.Error().Err(err).Int("order_id", order.ID).Msg("unexpected error while adding accrual to order")
		p.markOrderAsFailed(order)
	}
}

func (p *AccrualProcessor) markOrderAsFailed(order models.Order) {
	err := p.ordersStorage.ChangeStatus(order.ID, models.OrderStatusError)
	if err != nil {
		log.Error().
			Err(err).
			Int("order_id", order.ID).
			Msg("unexpected error while marking order as errored")
	}
}

func (p *AccrualProcessor) calculateAccrualForOrderItem(_ context.Context, orderItem models.OrderItem) (float64, error) {
	matchingReward, err := p.rewardsStorage.GetMatchingReward(orderItem)
	return matchingReward.CalculateReward(orderItem.Price), err
}
