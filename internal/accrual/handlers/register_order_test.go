package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/belamov/ypgo-gophermart/internal/accrual/mocks"
	"github.com/belamov/ypgo-gophermart/internal/accrual/models"
	"github.com/belamov/ypgo-gophermart/internal/accrual/services"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_RegisterOrder(t *testing.T) {
	validOrderID := "123"
	type want struct {
		statusCode int
	}
	tests := []struct {
		name    string
		want    want
		orderID string
		items   []models.OrderItem
		err     error
	}{
		{
			name: "it accepts new order",
			want: want{
				statusCode: http.StatusAccepted,
			},
			orderID: validOrderID,
			items:   []models.OrderItem{{Description: "item 1", Price: 700.5}},
		},
		{
			name: "it responds with 409 when order is already registered",
			want: want{
				statusCode: http.StatusConflict,
			},
			orderID: validOrderID,
			items:   []models.OrderItem{{Description: "item 1", Price: 700.5}},
			err:     services.ErrOrderIsAlreadyRegistered,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockOrderManager := mocks.NewMockOrderManagementInterface(ctrl)
			orderIDInt, err := strconv.Atoi(tt.orderID)
			require.NoError(t, err)
			if tt.err != nil {
				mockOrderManager.EXPECT().RegisterNewOrder(orderIDInt, tt.items).Return(tt.err).AnyTimes()
			} else {
				mockOrderManager.EXPECT().RegisterNewOrder(orderIDInt, tt.items).Return(nil).AnyTimes()
			}

			mockRewards := mocks.NewMockRewardsStorage(ctrl)

			r := NewRouter(mockOrderManager, mockRewards)
			ts := httptest.NewServer(r)
			defer ts.Close()

			request := newOrderRequest{
				Order: tt.orderID,
				Items: tt.items,
			}
			requestJSON, err := json.Marshal(request)
			require.NoError(t, err)
			result, _ := testRequest(
				t,
				ts,
				http.MethodPost,
				"/api/orders",
				string(requestJSON),
			)
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}
