package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/belamov/ypgo-gophermart/internal/accrual/mocks"
	"github.com/belamov/ypgo-gophermart/internal/accrual/models"
	"github.com/golang/mock/gomock"
)

func TestItProcessesOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderStorage := mocks.NewMockOrdersStorage(ctrl)
	mockOrderStorage.EXPECT().GetOrdersForProcessing().Return(nil, nil).Times(1)

	mockRewardStorage := mocks.NewMockRewardsStorage(ctrl)

	processor := NewAccrualProcessor(mockOrderStorage, mockRewardStorage)

	ctx, cancel := context.WithCancel(context.Background())

	go processor.StartProcessing(ctx)

	for i := 5; i < 9; i++ {
		item1 := models.OrderItem{Description: "", Price: float64(10 * i)}
		item2 := models.OrderItem{Description: "", Price: float64(50 * i)}
		item3 := models.OrderItem{Description: "", Price: float64(100 * i)}
		order := models.Order{
			ID:    i,
			Items: []models.OrderItem{item1, item2, item3},
		}

		absoluteReward := models.Reward{
			Match:      "",
			Reward:     10,
			RewardType: models.AbsoluteType,
		}
		percentReward := models.Reward{
			Match:      "",
			Reward:     10,
			RewardType: models.PercentType,
		}
		mockOrderStorage.EXPECT().ChangeStatus(order.ID, models.OrderStatusProcessing).Return(nil)
		mockRewardStorage.EXPECT().GetMatchingReward(item1).Return(absoluteReward, nil)
		mockRewardStorage.EXPECT().GetMatchingReward(item2).Return(percentReward, nil)
		mockRewardStorage.EXPECT().GetMatchingReward(item3).Return(models.Reward{}, nil)

		accrual := absoluteReward.CalculateReward(item1.Price) + percentReward.CalculateReward(item2.Price)
		mockOrderStorage.EXPECT().AddAccrual(order.ID, accrual).Return(nil)

		go processor.RegisterOrderForProcessing(order)
	}
	time.Sleep(time.Millisecond * 100)
	cancel()
}

func TestItProcessesErroredOrdersOnInit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderStorage := mocks.NewMockOrdersStorage(ctrl)

	mockRewardStorage := mocks.NewMockRewardsStorage(ctrl)

	processor := NewAccrualProcessor(mockOrderStorage, mockRewardStorage)

	ctx, cancel := context.WithCancel(context.Background())

	item1 := models.OrderItem{Description: "", Price: float64(10)}
	item2 := models.OrderItem{Description: "", Price: float64(20)}
	order := models.Order{
		ID:    1,
		Items: []models.OrderItem{item1, item2},
	}

	absoluteReward := models.Reward{
		Match:      "",
		Reward:     10,
		RewardType: models.AbsoluteType,
	}
	percentReward := models.Reward{
		Match:      "",
		Reward:     10,
		RewardType: models.PercentType,
	}
	mockOrderStorage.EXPECT().ChangeStatus(order.ID, models.OrderStatusProcessing).Return(nil)
	mockRewardStorage.EXPECT().GetMatchingReward(item1).Return(absoluteReward, nil)
	mockRewardStorage.EXPECT().GetMatchingReward(item2).Return(percentReward, nil)

	accrual := absoluteReward.CalculateReward(item1.Price) + percentReward.CalculateReward(item2.Price)
	mockOrderStorage.EXPECT().AddAccrual(order.ID, accrual).Return(nil)

	mockOrderStorage.EXPECT().GetOrdersForProcessing().Return([]models.Order{order}, nil).Times(1)

	go processor.StartProcessing(ctx)

	time.Sleep(time.Millisecond * 100)
	cancel()
}

func TestItMarksOrderAsErrored(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderStorage := mocks.NewMockOrdersStorage(ctrl)

	mockRewardStorage := mocks.NewMockRewardsStorage(ctrl)

	processor := NewAccrualProcessor(mockOrderStorage, mockRewardStorage)

	ctx, cancel := context.WithCancel(context.Background())

	item1 := models.OrderItem{Description: "", Price: float64(10)}
	item2 := models.OrderItem{Description: "", Price: float64(20)}
	order := models.Order{
		ID:    1,
		Items: []models.OrderItem{item1, item2},
	}

	percentReward := models.Reward{
		Match:      "",
		Reward:     10,
		RewardType: models.PercentType,
	}
	mockOrderStorage.EXPECT().ChangeStatus(order.ID, models.OrderStatusProcessing).Return(nil)
	mockOrderStorage.EXPECT().ChangeStatus(order.ID, models.OrderStatusError).Return(nil)
	mockRewardStorage.EXPECT().GetMatchingReward(item1).Return(models.Reward{}, errors.New("some"))
	mockRewardStorage.EXPECT().GetMatchingReward(item2).Return(percentReward, nil)

	processor.AddAccrualForOrder(ctx, order)

	time.Sleep(time.Millisecond * 100)
	cancel()
}
