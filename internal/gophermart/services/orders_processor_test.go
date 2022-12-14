package services

import (
	"context"
	"testing"
	"time"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/mocks"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/golang/mock/gomock"
)

func TestOrdersProcessor_ProcessOrderFirstFetchSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrdersStorage := mocks.NewMockOrdersStorage(ctrl)
	mockAccrual := mocks.NewMockAccrualInfoProvider(ctrl)
	mockBalance := mocks.NewMockBalanceProcessorInterface(ctrl)

	ordersProcessor := NewOrderProcessor(mockOrdersStorage, mockAccrual, mockBalance)

	order := models.Order{
		ID:         1,
		CreatedBy:  1,
		UploadedAt: time.Now(),
		Status:     models.OrderStatusNew,
		Accrual:    0,
	}

	mockOrdersStorage.EXPECT().ChangeStatus(order, models.OrderStatusProcessing).Return(nil).Times(1)
	mockAccrual.EXPECT().GetAccrualForOrder(gomock.Any(), order.ID).Return(100.0, nil).Times(1)
	mockBalance.EXPECT().AddAccrual(order, 100.0).Return(nil).Times(1)
	ordersProcessor.ProcessOrder(context.Background(), order)
}

func TestOrdersProcessor_ProcessOrderThirdFetchSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrdersStorage := mocks.NewMockOrdersStorage(ctrl)
	mockAccrual := mocks.NewMockAccrualInfoProvider(ctrl)
	mockBalance := mocks.NewMockBalanceProcessorInterface(ctrl)

	ordersProcessor := NewOrderProcessor(mockOrdersStorage, mockAccrual, mockBalance)

	order := models.Order{
		ID:         1,
		CreatedBy:  1,
		UploadedAt: time.Now(),
		Status:     models.OrderStatusNew,
		Accrual:    0,
	}

	mockOrdersStorage.EXPECT().ChangeStatus(order, models.OrderStatusProcessing).Return(nil).Times(1)
	mockAccrual.EXPECT().GetAccrualForOrder(gomock.Any(), order.ID).Return(0.0, ErrOrderIsNotYetProceeded).Times(1)
	mockAccrual.EXPECT().GetAccrualForOrder(gomock.Any(), order.ID).Return(0.0, ErrOrderIsNotYetProceeded).Times(1)
	mockAccrual.EXPECT().GetAccrualForOrder(gomock.Any(), order.ID).Return(100.0, nil).Times(1)
	mockBalance.EXPECT().AddAccrual(order, 100.0).Return(nil).Times(1)
	ordersProcessor.ProcessOrder(context.Background(), order)
}

func TestOrdersProcessor_ProcessOrderInvalidAccrual(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrdersStorage := mocks.NewMockOrdersStorage(ctrl)
	mockAccrual := mocks.NewMockAccrualInfoProvider(ctrl)
	mockBalance := mocks.NewMockBalanceProcessorInterface(ctrl)

	ordersProcessor := NewOrderProcessor(mockOrdersStorage, mockAccrual, mockBalance)

	order := models.Order{
		ID:         1,
		CreatedBy:  1,
		UploadedAt: time.Now(),
		Status:     models.OrderStatusNew,
		Accrual:    0,
	}

	mockOrdersStorage.EXPECT().ChangeStatus(order, models.OrderStatusProcessing).Return(nil).Times(1)
	mockAccrual.EXPECT().GetAccrualForOrder(gomock.Any(), order.ID).Return(0.0, ErrOrderIsNotYetProceeded).Times(1)
	mockAccrual.EXPECT().GetAccrualForOrder(gomock.Any(), order.ID).Return(0.0, ErrInvalidOrderForAccrual).Times(1)
	mockOrdersStorage.EXPECT().ChangeStatus(order, models.OrderStatusInvalid).Return(nil).Times(1)
	ordersProcessor.ProcessOrder(context.Background(), order)
}

func TestOrdersProcessor_ItProcessingOrdersInBackground(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrdersStorage := mocks.NewMockOrdersStorage(ctrl)
	mockAccrual := mocks.NewMockAccrualInfoProvider(ctrl)
	mockBalance := mocks.NewMockBalanceProcessorInterface(ctrl)

	ordersProcessor := NewOrderProcessor(mockOrdersStorage, mockAccrual, mockBalance)

	ctx, cancel := context.WithCancel(context.Background())

	mockOrdersStorage.EXPECT().GetOrdersForProcessing().Return(nil, nil).Times(1)
	mockOrdersStorage.EXPECT().ChangeStatus(gomock.Any(), models.OrderStatusProcessing).Return(nil).MinTimes(1)
	mockAccrual.EXPECT().GetAccrualForOrder(ctx, gomock.Any()).Return(0.0, ErrOrderIsNotYetProceeded).MinTimes(1)

	go ordersProcessor.StartProcessing(ctx)

	for i := 0; i < 100; i++ {
		order := models.Order{
			ID:         i + 1,
			CreatedBy:  1,
			UploadedAt: time.Now(),
			Status:     models.OrderStatusNew,
			Accrual:    0,
		}
		go ordersProcessor.RegisterOrderForProcessing(order)
	}

	mockOrdersStorage.EXPECT().ChangeStatus(gomock.Any(), models.OrderStatusNew).Return(nil).MinTimes(1)
	time.Sleep(time.Millisecond * 500)
	cancel()
	time.Sleep(time.Millisecond * 500)
}

func TestOrdersProcessor_ItStartsProcessingOrdersThatMustBeProcessed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrdersStorage := mocks.NewMockOrdersStorage(ctrl)
	mockAccrual := mocks.NewMockAccrualInfoProvider(ctrl)
	mockBalance := mocks.NewMockBalanceProcessorInterface(ctrl)

	ordersProcessor := NewOrderProcessor(mockOrdersStorage, mockAccrual, mockBalance)

	ctx, cancel := context.WithCancel(context.Background())

	newOrders := []models.Order{
		{ID: 1, CreatedBy: 1, UploadedAt: time.Now(), Status: models.OrderStatusNew, Accrual: 0},
		{ID: 2, CreatedBy: 1, UploadedAt: time.Now(), Status: models.OrderStatusNew, Accrual: 0},
	}

	mockOrdersStorage.EXPECT().GetOrdersForProcessing().Return(newOrders, nil).Times(1)

	for _, newOrder := range newOrders {
		mockOrdersStorage.EXPECT().ChangeStatus(newOrder, gomock.Any()).Return(nil).MinTimes(1)
		mockAccrual.EXPECT().GetAccrualForOrder(ctx, newOrder.ID).Return(0.0, ErrOrderIsNotYetProceeded).MinTimes(1)
	}

	go ordersProcessor.StartProcessing(ctx)
	time.Sleep(time.Millisecond * 5)

	cancel()
}

func TestOrdersProcessor_ItProcessesOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrdersStorage := mocks.NewMockOrdersStorage(ctrl)
	mockAccrual := mocks.NewMockAccrualInfoProvider(ctrl)
	mockBalance := mocks.NewMockBalanceProcessorInterface(ctrl)

	ordersProcessor := NewOrderProcessor(mockOrdersStorage, mockAccrual, mockBalance)

	ctx, cancel := context.WithCancel(context.Background())

	newOrders := []models.Order{
		{ID: 1, CreatedBy: 1, UploadedAt: time.Now(), Status: models.OrderStatusNew, Accrual: 0},
		{ID: 2, CreatedBy: 1, UploadedAt: time.Now(), Status: models.OrderStatusNew, Accrual: 0},
	}

	mockOrdersStorage.EXPECT().GetOrdersForProcessing().Return(nil, nil).Times(1)

	for _, newOrder := range newOrders {
		mockOrdersStorage.EXPECT().ChangeStatus(newOrder, gomock.Any()).Return(nil).Times(1)
		mockAccrual.EXPECT().GetAccrualForOrder(ctx, newOrder.ID).Return(5.0, nil).Times(1)
		mockBalance.EXPECT().AddAccrual(newOrder, 5.0).Return(nil).Times(1)
	}

	go ordersProcessor.StartProcessing(ctx)

	for _, newOrder := range newOrders {
		go ordersProcessor.RegisterOrderForProcessing(newOrder)
	}
	time.Sleep(time.Millisecond * 5)
	cancel()
}
