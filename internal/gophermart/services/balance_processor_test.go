package services

import (
	"testing"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/mocks"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestBalanceProcessor_RegisterWithdraw(t *testing.T) {
	userWithSufficientBalance := models.User{ID: 1}
	userWithInsufficientBalance := models.User{ID: 2}
	tests := []struct {
		name    string
		orderID int
		userID  int
		amount  float64
		wantErr bool
	}{
		{
			name:    "it registers withdraw",
			orderID: 1,
			userID:  userWithSufficientBalance.ID,
			amount:  100,
			wantErr: false,
		},
		{
			name:    "it doesn't register withdraw when user has no balance",
			orderID: 1,
			userID:  userWithInsufficientBalance.ID,
			amount:  100,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockBalanceStorage := mocks.NewMockBalanceStorage(ctrl)
			mockBalanceStorage.EXPECT().GetTotalAccrual(userWithSufficientBalance.ID).Return(100000.0, nil).AnyTimes()
			mockBalanceStorage.EXPECT().GetTotalAccrual(userWithInsufficientBalance.ID).Return(0.0, nil).AnyTimes()
			mockBalanceStorage.EXPECT().GetTotalWithdraws(userWithSufficientBalance.ID).Return(0.0, nil).AnyTimes()
			mockBalanceStorage.EXPECT().GetTotalWithdraws(userWithInsufficientBalance.ID).Return(0.0, nil).AnyTimes()
			mockBalanceStorage.EXPECT().AddWithdraw(gomock.Any(), userWithSufficientBalance.ID, gomock.Any()).Return(nil).AnyTimes()

			b := NewBalanceProcessor(mockBalanceStorage)
			err := b.RegisterWithdraw(tt.orderID, tt.userID, tt.amount)
			if tt.wantErr {
				assert.Error(t, err)
				var ibe *InsufficientBalanceError
				assert.ErrorAs(t, err, &ibe)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
