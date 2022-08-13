package services

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccrualHttpClient_GetAccrualForOrderSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		_, err := rw.Write([]byte(`
				{
					"order": "100",
					"status": "PROCESSED",
					"accrual": 50.3
				}
		`))
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewAccrualHTTPClient(server.Client(), server.URL, 50)
	accrual, err := client.GetAccrualForOrder(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, 50.3, accrual)
}

func TestAccrualHttpClient_GetAccrualForOrderInvalid(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		_, err := rw.Write([]byte(`
				{
					"order": "100",
					"status": "INVALID"
				}
		`))
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewAccrualHTTPClient(server.Client(), server.URL, 50)
	accrual, err := client.GetAccrualForOrder(context.Background(), 1)
	assert.ErrorIs(t, err, ErrInvalidOrderForAccrual)
	assert.Equal(t, 0.0, accrual)
}

func TestAccrualHttpClient_GetAccrualForOrderProcessing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		_, err := rw.Write([]byte(`
				{
					"order": "100",
					"status": "PROCESSING"
				}
		`))
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewAccrualHTTPClient(server.Client(), server.URL, 50)
	accrual, err := client.GetAccrualForOrder(context.Background(), 1)
	assert.ErrorIs(t, err, ErrOrderIsNotYetProceeded)
	assert.Equal(t, 0.0, accrual)
}

func TestAccrualHttpClient_GetAccrualForOrderRegistered(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		_, err := rw.Write([]byte(`
				{
					"order": "100",
					"status": "REGISTERED"
				}
		`))
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewAccrualHTTPClient(server.Client(), server.URL, 50)
	accrual, err := client.GetAccrualForOrder(context.Background(), 1)
	assert.ErrorIs(t, err, ErrOrderIsNotYetProceeded)
	assert.Equal(t, 0.0, accrual)
}
