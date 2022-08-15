package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/mocks"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/services"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_RegisterWithdraw(t *testing.T) {
	validOrderID := "123"
	invalidOrderID := "456"
	user := models.User{
		ID:             10,
		Login:          "login",
		HashedPassword: "hash",
	}
	userWithInsufficientBalance := models.User{ID: 20}
	type want struct {
		statusCode int
	}
	tests := []struct {
		name    string
		want    want
		orderID string
		amount  float64
		userID  int
	}{
		{
			name: "it accepts valid order id and valid amount",
			want: want{
				statusCode: http.StatusOK,
			},
			orderID: validOrderID,
			amount:  100.5,
		},
		{
			name: "it doesnt accept invalid order id and valid amount",
			want: want{
				statusCode: http.StatusUnprocessableEntity,
			},
			orderID: invalidOrderID,
			amount:  100.5,
		},
		{
			name: "it doesnt accept valid order id and invalid amount (zero)",
			want: want{
				statusCode: http.StatusUnprocessableEntity,
			},
			orderID: validOrderID,
			amount:  0,
		},
		{
			name: "it doesnt accept valid order id and invalid amount (negative)",
			want: want{
				statusCode: http.StatusUnprocessableEntity,
			},
			orderID: validOrderID,
			amount:  -100,
		},
		{
			name: "it doesnt accept invalid order id (letters)",
			want: want{
				statusCode: http.StatusUnprocessableEntity,
			},
			orderID: "e1a4",
			amount:  100,
		},
		{
			name: "it doesnt accept invalid order id (empty)",
			want: want{
				statusCode: http.StatusUnprocessableEntity,
			},
			orderID: "",
			amount:  100,
		},
		{
			name: "it returns 402 when user has insufficient balance",
			want: want{
				statusCode: http.StatusPaymentRequired,
			},
			orderID: validOrderID,
			amount:  100,
			userID:  userWithInsufficientBalance.ID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var mockAuth *mocks.MockAuth
			if tt.userID == 0 {
				mockAuth = mocks.NewMockAuth(ctrl)
				mockAuth.EXPECT().AuthMiddleware().Return(emptyMiddleware).AnyTimes()
				mockAuth.EXPECT().GetUserID(gomock.Any()).Return(user.ID).AnyTimes()
			} else {
				mockAuth = mocks.NewMockAuth(ctrl)
				mockAuth.EXPECT().AuthMiddleware().Return(emptyMiddleware).AnyTimes()
				mockAuth.EXPECT().GetUserID(gomock.Any()).Return(userWithInsufficientBalance.ID).AnyTimes()
			}

			mockOrders := mocks.NewMockOrdersManagerInterface(ctrl)
			validOrderIDInt, err := strconv.Atoi(validOrderID)
			require.NoError(t, err)
			invalidOrderIDInt, err := strconv.Atoi(invalidOrderID)
			require.NoError(t, err)
			mockOrders.EXPECT().ValidateOrderID(invalidOrderIDInt).Return(errors.New("invalid order number")).AnyTimes()
			mockOrders.EXPECT().ValidateOrderID(validOrderIDInt).Return(nil).AnyTimes()

			mockBalance := mocks.NewMockBalanceProcessorInterface(ctrl)
			mockBalance.EXPECT().RegisterWithdraw(validOrderIDInt, user.ID, tt.amount).Return(nil).AnyTimes()
			mockBalance.EXPECT().RegisterWithdraw(validOrderIDInt, userWithInsufficientBalance.ID, tt.amount).Return(services.NewInsufficientBalanceError(0, tt.amount)).AnyTimes()

			r := NewRouter(mockAuth, mockOrders, mockBalance)
			ts := httptest.NewServer(r)
			defer ts.Close()

			request := RegisterWithdrawRequest{
				OrderID: tt.orderID,
				Amount:  tt.amount,
			}
			requestJSON, err := json.Marshal(request)
			require.NoError(t, err)

			result, _ := testRequest(
				t,
				ts,
				http.MethodPost,
				"/api/user/balance/withdraw",
				string(requestJSON),
			)
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}
