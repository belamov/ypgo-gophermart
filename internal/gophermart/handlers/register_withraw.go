package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/services"
)

type RegisterWithdrawRequest struct {
	OrderID string  `json:"order"`
	Amount  float64 `json:"sum"`
}

func (h *Handler) RegisterWithdraw(w http.ResponseWriter, r *http.Request) {
	reader, err := getDecompressedReader(r)
	if err != nil {
		log.Println("unexpected error in register withdraws handler:")
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var request RegisterWithdrawRequest
	if err := json.NewDecoder(reader).Decode(&request); err != nil {
		log.Println("unexpected error in register withdraws handler:")
		log.Println(err.Error())
		http.Error(w, "cannot decode json", http.StatusBadRequest)
		return
	}

	if request.Amount <= 0 || request.OrderID == "" {
		http.Error(w, "sum must be greater than 0, order id is required", http.StatusUnprocessableEntity)
		return
	}

	userID := h.auth.GetUserID(r)
	if userID == 0 {
		log.Println("unexpected error in register withdraws handler:")
		log.Println("user id not found")
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
	var insufficientBalanceError *services.InsufficientBalanceError
	if errors.As(err, &insufficientBalanceError) {
		http.Error(w, err.Error(), http.StatusPaymentRequired)
		return
	}
	if err != nil {
		log.Println("unexpected error in register withdraws handler:")
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
