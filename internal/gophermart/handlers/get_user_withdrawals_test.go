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

func TestHandler_GetUserWithdrawals(t *testing.T) {
	type want struct {
		json   string
		status int
	}
	tests := []struct {
		name        string
		withdrawals []models.Withdrawal
		want        want
	}{
		{
			name:        "it responds with 204 when list is empty",
			withdrawals: []models.Withdrawal{},
			want: want{
				status: http.StatusNoContent,
				json:   "",
			},
		},
		{
			name: "it responds with 200 when list is not empty",
			withdrawals: []models.Withdrawal{
				{
					OrderID:          9278923470,
					CreatedAt:        getTimeFromString("2020-12-10T15:15:45+03:00"),
					WithdrawalAmount: 500.5,
				},
				{
					OrderID:          346436439,
					CreatedAt:        getTimeFromString("2020-12-09T16:09:53+03:00"),
					WithdrawalAmount: 10,
				},
			},
			want: want{
				status: http.StatusOK,
				json: `[
					{
						"order": "9278923470",
						"sum": 500.5,
						"processed_at": "2020-12-10T15:15:45+03:00"
					},
					{
						"order": "346436439",
						"sum": 10,
						"processed_at": "2020-12-09T16:09:53+03:00"
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

			mockBalance := mocks.NewMockBalanceProcessorInterface(ctrl)
			mockBalance.EXPECT().GetUserWithdrawals(gomock.Any()).Return(tt.withdrawals, nil).AnyTimes()

			r := NewRouter(mockAuth, mockOrders, mockBalance)
			ts := httptest.NewServer(r)
			defer ts.Close()

			result, body := testRequest(
				t,
				ts,
				http.MethodGet,
				"/api/user/withdrawals",
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
