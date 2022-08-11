package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_GetUserBalance(t *testing.T) {
	type want struct {
		json   string
		status int
	}
	tests := []struct {
		name      string
		withdrawn float64
		accrual   float64
		want      want
	}{
		{
			name:      "it returns user balance",
			withdrawn: 42,
			accrual:   542.5,
			want: want{
				status: http.StatusOK,
				json: `
					{
						"current": 500.5,
						"withdrawn": 42
					}
				`,
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
			mockBalance.EXPECT().GetUserTotalWithdrawAmount(gomock.Any()).Return(tt.withdrawn, nil).AnyTimes()
			mockBalance.EXPECT().GetUserTotalAccrualAmount(gomock.Any()).Return(tt.accrual, nil).AnyTimes()

			r := NewRouter(mockAuth, mockOrders, mockBalance)
			ts := httptest.NewServer(r)
			defer ts.Close()

			result, body := testRequest(
				t,
				ts,
				http.MethodGet,
				"/api/user/balance",
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
