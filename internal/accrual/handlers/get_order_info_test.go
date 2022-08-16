package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/belamov/ypgo-gophermart/internal/accrual/mocks"
	"github.com/belamov/ypgo-gophermart/internal/accrual/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_GetOrderInfoRegistered(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	order := models.Order{
		ID:      22,
		Status:  models.OrderStatusNew,
		Accrual: 15.2,
	}
	mockOrderManager := mocks.NewMockOrderManagementInterface(ctrl)

	mockOrderManager.EXPECT().ValidateOrderID(order.ID).Return(nil).AnyTimes()
	mockOrderManager.EXPECT().GetOrderInfo(order.ID).Return(order, nil)
	mockRewards := mocks.NewMockRewardsStorage(ctrl)

	r := NewRouter(mockOrderManager, mockRewards)
	ts := httptest.NewServer(r)
	defer ts.Close()

	result, body := testRequest(
		t,
		ts,
		http.MethodGet,
		"/api/orders/"+strconv.Itoa(order.ID),
		"",
	)
	defer result.Body.Close()

	response := OrderResponse{
		Number:  strconv.Itoa(order.ID),
		Status:  order.Status.String(),
		Accrual: order.Accrual,
	}

	out, err := json.Marshal(response)
	require.NoError(t, err)
	assert.JSONEq(t, string(out), body)
}

func TestHandler_GetOrderInfoInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	order := models.Order{
		ID: 22,
	}
	mockOrderManager := mocks.NewMockOrderManagementInterface(ctrl)

	mockOrderManager.EXPECT().ValidateOrderID(order.ID).Return(errors.New("")).AnyTimes()
	mockRewards := mocks.NewMockRewardsStorage(ctrl)

	r := NewRouter(mockOrderManager, mockRewards)
	ts := httptest.NewServer(r)
	defer ts.Close()

	result, body := testRequest(
		t,
		ts,
		http.MethodGet,
		"/api/orders/"+strconv.Itoa(order.ID),
		"",
	)
	defer result.Body.Close()

	response := OrderResponse{
		Number:  strconv.Itoa(order.ID),
		Status:  models.OrderStatusInvalid.String(),
		Accrual: order.Accrual,
	}

	out, err := json.Marshal(response)
	require.NoError(t, err)
	assert.JSONEq(t, string(out), body)
}
