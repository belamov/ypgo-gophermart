package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/services"
	"github.com/rs/zerolog/log"
)

type RegisterWithdrawRequest struct {
	OrderID string  `json:"order"`
	Amount  float64 `json:"sum"`
}

func (h *Handler) RegisterWithdraw(w http.ResponseWriter, r *http.Request) {
	reader, err := getDecompressedReader(r)
	if err != nil {
		log.Error().Err(err).Msg("unexpected error in register withdraws handler")
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
		log.Error().Err(err).Msg("unexpected error in register withdraws handler. user id not found")
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	orderID, err := strconv.Atoi(request.OrderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err = h.ordersManager.ValidateOrderID(orderID); err != nil {
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
		log.Error().Err(err).Msg("unexpected error in register withdraws handler")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
