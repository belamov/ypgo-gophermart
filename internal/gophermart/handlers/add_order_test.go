package handlers

import (
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

func TestHandler_AddOrder(t *testing.T) {
	validOrderID := "123"
	invalidOrderID := "456"
	user := models.User{
		ID:             10,
		Login:          "login",
		HashedPassword: "hash",
	}
	type want struct {
		statusCode int
	}
	tests := []struct {
		name    string
		orderID string
		user    models.User
		want    want
	}{
		{
			name: "it accepts valid order id from authenticated user",
			want: want{
				statusCode: http.StatusAccepted,
			},
			orderID: validOrderID,
			user:    user,
		},
		{
			name: "it doesnt accept valid order id from unauthenticated user",
			want: want{
				statusCode: http.StatusUnauthorized,
			},
			orderID: validOrderID,
			user:    models.User{},
		},
		{
			name: "it doesnt accept invalid order id from authenticated user",
			want: want{
				statusCode: http.StatusUnprocessableEntity,
			},
			orderID: invalidOrderID,
			user:    user,
		},
		{
			name: "it doesnt accept invalid order id (order id has letters)",
			want: want{
				statusCode: http.StatusUnprocessableEntity,
			},
			orderID: "12e23a",
			user:    user,
		},
		{
			name: "it doesnt accept invalid order id (order id has only letters)",
			want: want{
				statusCode: http.StatusUnprocessableEntity,
			},
			orderID: "some id",
			user:    user,
		},
		{
			name: "it responds with 200 when order already added by same user",
			want: want{
				statusCode: http.StatusOK,
			},
			orderID: "111",
			user:    user,
		},
		{
			name: "it responds with 409 when order already added by another user",
			want: want{
				statusCode: http.StatusConflict,
			},
			orderID: "222",
			user:    user,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuth := mocks.NewMockAuth(ctrl)
			mockAuth.EXPECT().AuthMiddleware().Return(emptyMiddleware).AnyTimes()
			mockAuth.EXPECT().GetUserID(gomock.Any()).Return(tt.user.ID).AnyTimes()

			mockOrders := mocks.NewMockOrdersManagerInterface(ctrl)

			validOrderIDInt, err := strconv.Atoi(validOrderID)
			require.NoError(t, err)
			invalidOrderIDInt, err := strconv.Atoi(invalidOrderID)
			require.NoError(t, err)
			mockOrders.EXPECT().AddOrder(validOrderIDInt, tt.user.ID).Return(nil).AnyTimes()
			mockOrders.EXPECT().AddOrder(111, tt.user.ID).Return(services.NewOrderAlreadyAddedError(models.Order{ID: 111, CreatedBy: tt.user.ID})).AnyTimes()
			mockOrders.EXPECT().AddOrder(222, tt.user.ID).Return(services.NewOrderAlreadyAddedError(models.Order{ID: 111, CreatedBy: 20})).AnyTimes()
			mockOrders.EXPECT().ValidateOrderID(invalidOrderIDInt).Return(errors.New("invalid order number")).AnyTimes()
			mockOrders.EXPECT().ValidateOrderID(validOrderIDInt).Return(nil).AnyTimes()
			mockOrders.EXPECT().ValidateOrderID(111).Return(nil).AnyTimes()
			mockOrders.EXPECT().ValidateOrderID(222).Return(nil).AnyTimes()

			mockBalance := mocks.NewMockBalanceProcessorInterface(ctrl)

			r := NewRouter(mockAuth, mockOrders, mockBalance)
			ts := httptest.NewServer(r)
			defer ts.Close()

			result, _ := testRequest(
				t,
				ts,
				http.MethodPost,
				"/api/user/orders",
				tt.orderID,
			)
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}
