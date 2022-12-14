package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/belamov/ypgo-gophermart/internal/accrual/models"
	"github.com/belamov/ypgo-gophermart/internal/accrual/services"
	"github.com/rs/zerolog/log"
)

type newOrderRequest struct {
	Order string             `json:"order"`
	Items []models.OrderItem `json:"goods"`
}

func (h *Handler) RegisterOrder(w http.ResponseWriter, r *http.Request) {
	reader, err := getDecompressedReader(r)
	if err != nil {
		log.Error().Err(err).Msg("unexpected error in register order handler:")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var newOrder newOrderRequest

	if err = json.NewDecoder(reader).Decode(&newOrder); err != nil {
		http.Error(w, "cannot decode json", http.StatusBadRequest)
		return
	}

	orderID, err := strconv.Atoi(newOrder.Order)
	if err != nil {
		http.Error(w, "wrong order id", http.StatusBadRequest)
		return
	}

	if err = h.orderManager.ValidateOrderID(orderID); err != nil {
		http.Error(w, "wrong order id", http.StatusBadRequest)
		return
	}

	err = h.orderManager.RegisterNewOrder(orderID, newOrder.Items)
	if errors.Is(err, services.ErrOrderIsAlreadyRegistered) {
		w.WriteHeader(http.StatusConflict)
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("unexpected error in register order handler:")
		http.Error(w, "wrong order id", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
