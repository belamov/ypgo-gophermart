package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/mocks"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_GetUsersOrders(t *testing.T) {
	type want struct {
		json   string
		status int
	}
	tests := []struct {
		name   string
		orders []models.Order
		want   want
	}{
		{
			name:   "it responds with 204 when list is empty",
			orders: []models.Order{},
			want: want{
				status: http.StatusNoContent,
				json:   "",
			},
		},
		{
			name: "it responds with 200 when list is not empty",
			orders: []models.Order{
				{
					ID:         9278923470,
					UploadedAt: getTimeFromString("2020-12-10T15:15:45+03:00"),
					Status:     models.OrderStatusProcessed,
					Accrual:    500,
				},
				{
					ID:         12345678903,
					UploadedAt: getTimeFromString("2020-12-10T15:12:01+03:00"),
					Status:     models.OrderStatusProcessing,
					Accrual:    0,
				},
				{
					ID:         346436439,
					UploadedAt: getTimeFromString("2020-12-09T16:09:53+03:00"),
					Status:     models.OrderStatusInvalid,
					Accrual:    0,
				},
			},
			want: want{
				status: http.StatusOK,
				json: `[
					{
						"number": "9278923470",
						"status": "PROCESSED",
						"accrual": 500,
						"uploaded_at": "2020-12-10T15:15:45+03:00"
					},
					{
						"number": "12345678903",
						"status": "PROCESSING",
						"uploaded_at": "2020-12-10T15:12:01+03:00"
					},
					{
						"number": "346436439",
						"status": "INVALID",
						"uploaded_at": "2020-12-09T16:09:53+03:00"
					}
				]`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuth := mocks.NewMockAuth(ctrl)
			mockAuth.EXPECT().AuthMiddleware().Return(emptyMiddleware).AnyTimes()
			mockAuth.EXPECT().GetUserID(gomock.Any()).Return(1).AnyTimes()

			mockOrders := mocks.NewMockOrdersProcessorInterface(ctrl)
			mockOrders.EXPECT().GetUsersOrders(gomock.Any()).Return(tt.orders, nil).AnyTimes()

			mockBalance := mocks.NewMockBalanceProcessorInterface(ctrl)

			r := NewRouter(mockAuth, mockOrders, mockBalance)
			ts := httptest.NewServer(r)
			defer ts.Close()

			result, body := testRequest(
				t,
				ts,
				http.MethodGet,
				"/api/user/orders",
				"",
			)
			defer result.Body.Close()

			assert.Equal(t, tt.want.status, result.StatusCode)
			if tt.want.json != "" {
				assert.JSONEq(t, tt.want.json, body)
			}
		})
	}
}
