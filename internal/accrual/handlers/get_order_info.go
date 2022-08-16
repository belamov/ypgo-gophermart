package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/belamov/ypgo-gophermart/internal/accrual/models"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type OrderResponse struct {
	Number  string  `json:"number"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}

func (h *Handler) GetOrderInfo(w http.ResponseWriter, r *http.Request) {
	orderIDRaw := chi.URLParam(r, "order")
	orderID, err := strconv.Atoi(orderIDRaw)
	if err != nil {
		http.Error(w, "wrong order id", http.StatusBadRequest)
		return
	}

	if err = h.orderManager.ValidateOrderID(orderID); err != nil {
		invalidResponse(w, orderID)
		return
	}

	order, err := h.orderManager.GetOrderInfo(orderID)
	if err != nil {
		log.Error().Err(err).Msg("unexpected error in get order handler")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if order.ID == 0 {
		invalidResponse(w, orderID)
	}

	w.Header().Set("Content-Type", "application/json")
	writeResponse(w, order)
}

func invalidResponse(w http.ResponseWriter, orderID int) {
	order := models.Order{
		ID:      orderID,
		Status:  models.OrderStatusInvalid,
		Accrual: 0,
	}
	writeResponse(w, order)
}

func writeResponse(w http.ResponseWriter, order models.Order) {
	response := OrderResponse{
		Number:  strconv.Itoa(order.ID),
		Status:  order.Status.String(),
		Accrual: order.Accrual,
	}

	out, err := json.Marshal(response)
	if err != nil {
		log.Error().Err(err).Msg("unexpected error in get order handler")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(out); err != nil {
		log.Error().Err(err).Msg("unexpected error in get order handler")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
