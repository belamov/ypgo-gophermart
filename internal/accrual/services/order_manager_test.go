package services

import (
	"testing"

	"github.com/belamov/ypgo-gophermart/internal/accrual/mocks"
	"github.com/belamov/ypgo-gophermart/internal/accrual/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestOrderManager_RegisterNewOrder(t *testing.T) {
	alreadyRegisteredOrderID := 1
	newOrderID := 2

	tests := []struct {
		name       string
		orderID    int
		orderItems []models.OrderItem
		wantErr    error
	}{
		{name: "it registering new order", orderID: newOrderID, orderItems: nil, wantErr: nil},
		{name: "it doesnt register existing order", orderID: alreadyRegisteredOrderID, orderItems: nil, wantErr: ErrOrderIsAlreadyRegistered},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := mocks.NewMockOrdersStorage(ctrl)
			mockStorage.EXPECT().Exists(alreadyRegisteredOrderID).Return(true, nil).AnyTimes()
			mockStorage.EXPECT().Exists(newOrderID).Return(false, nil).AnyTimes()
			mockStorage.EXPECT().CreateNew(tt.orderID, tt.orderItems).Return(nil).AnyTimes()

			service := NewOrderManager(mockStorage)

			err := service.RegisterNewOrder(tt.orderID, tt.orderItems)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
