package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type RegisterWithdrawRequest struct {
	OrderID string  `json:"order"`
	Amount  float64 `json:"sum"`
}

func (h *Handler) RegisterWithdraw(w http.ResponseWriter, r *http.Request) {
	reader, err := getDecompressedReader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var request RegisterWithdrawRequest
	if err := json.NewDecoder(reader).Decode(&request); err != nil {
		http.Error(w, "cannot decode json", http.StatusBadRequest)
		return
	}

	if request.Amount <= 0 || request.OrderID == "" {
		http.Error(w, "sum must be greater than 0, order id is required", http.StatusUnprocessableEntity)
		return
	}

	userID := h.auth.GetUserID(r)
	if userID == 0 {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	orderID, err := strconv.Atoi(request.OrderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err = h.ordersProcessor.ValidateOrderID(orderID); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err = h.balanceProcessor.RegisterWithdraw(orderID, userID, request.Amount)
	// TODO: handle insufficient balance error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
