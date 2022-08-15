package services

import (
	"errors"
	"testing"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/mocks"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestOrdersManager_ValidateOrderID(t *testing.T) {
	tests := []struct {
		name    string
		orderID int
		wantErr bool
	}{
		{name: "negative number", orderID: -100, wantErr: true},
		{name: "0 number", orderID: 0, wantErr: true},
		{name: "invalid by luhn", orderID: 4561261212345464, wantErr: true},
		{name: "valid by luhn", orderID: 4561261212345467, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockOrdersStorage := mocks.NewMockOrdersStorage(ctrl)
			mockAccrual := mocks.NewMockAccrualInfoProvider(ctrl)
			mockBalance := mocks.NewMockBalanceProcessorInterface(ctrl)

			ordersManager := NewOrdersManager(mockOrdersStorage, mockBalance, NewOrderProcessor(mockOrdersStorage, mockAccrual, mockBalance))
			err := ordersManager.ValidateOrderID(tt.orderID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestOrdersManager_AddOrder(t *testing.T) {
	tests := []struct {
		name    string
		orderID int
		userID  int
		wantErr bool
	}{
		{name: "existing order, existing user", orderID: 1, userID: 1, wantErr: true},
		{name: "existing order, new user", orderID: 1, userID: 2, wantErr: true},
		{name: "new order, new user", orderID: 2, userID: 2, wantErr: true},
		{name: "new order, existing user", orderID: 2, userID: 1, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockOrdersStorage := mocks.NewMockOrdersStorage(ctrl)
			mockOrdersStorage.EXPECT().FindByID(1).Return(models.Order{ID: 1, CreatedBy: 1}, nil).AnyTimes()
			mockOrdersStorage.EXPECT().FindByID(2).Return(models.Order{}, nil).AnyTimes()
			mockOrdersStorage.EXPECT().CreateNew(2, 1).Return(models.Order{ID: 1, CreatedBy: 1}, nil).AnyTimes()
			mockOrdersStorage.EXPECT().CreateNew(2, 2).Return(models.Order{}, errors.New("user id not found")).AnyTimes()

			mockAccrual := mocks.NewMockAccrualInfoProvider(ctrl)
			mockBalance := mocks.NewMockBalanceProcessorInterface(ctrl)

			mockOrdersStorage.EXPECT().ChangeStatus(gomock.Any(), models.OrderStatusProcessing).Return(nil).AnyTimes()
			mockAccrual.EXPECT().GetAccrualForOrder(gomock.Any(), gomock.Any()).Return(100.0, nil).AnyTimes()
			mockBalance.EXPECT().AddAccrual(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

			ordersProcessor := NewOrdersManager(mockOrdersStorage, mockBalance, &OrderProcessor{})
			err := ordersProcessor.AddOrder(tt.orderID, tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
