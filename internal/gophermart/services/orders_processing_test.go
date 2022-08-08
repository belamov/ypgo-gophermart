package services

import (
	"errors"
	"testing"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/mocks"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestOrdersProcessor_ValidateOrderID(t *testing.T) {
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

			ordersProcessor := NewOrdersProcessor(mockOrdersStorage)
			err := ordersProcessor.ValidateOrderID(tt.orderID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestOrdersProcessor_AddOrder(t *testing.T) {
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
			mockOrdersStorage.EXPECT().FindByID(1).Return(models.Order{ID: 1, UserID: 1}, nil).AnyTimes()
			mockOrdersStorage.EXPECT().FindByID(2).Return(models.Order{}, nil).AnyTimes()
			mockOrdersStorage.EXPECT().CreateNew(2, 1).Return(models.Order{ID: 1, UserID: 1}, nil).AnyTimes()
			mockOrdersStorage.EXPECT().CreateNew(2, 2).Return(models.Order{}, errors.New("user id not found")).AnyTimes()

			ordersProcessor := NewOrdersProcessor(mockOrdersStorage)
			err := ordersProcessor.AddOrder(tt.orderID, tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
