package services

import (
	"testing"

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
			ordersProcessor := NewOrdersProcessor()
			err := ordersProcessor.ValidateOrderID(tt.orderID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
