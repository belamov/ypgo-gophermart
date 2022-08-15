package handlers

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/services"
	"github.com/rs/zerolog/log"
)

func (h *Handler) AddOrder(w http.ResponseWriter, r *http.Request) {
	reader, err := getDecompressedReader(r)
	if err != nil {
		log.Error().Err(err).Msg("unexpected error in add order handler")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rawOrderID, err := io.ReadAll(reader)
	if err != nil {
		log.Error().Err(err).Msg("unexpected error in add order handler")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	orderID, err := strconv.Atoi(string(rawOrderID))
	if err != nil {
		http.Error(w, "invalid order id", http.StatusUnprocessableEntity)
		return
	}

	if err = h.ordersManager.ValidateOrderID(orderID); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	userID := h.auth.GetUserID(r)
	if userID == 0 {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	err = h.ordersManager.AddOrder(orderID, userID)
	if err != nil {
		handleOrderAddError(err, userID, w)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func handleOrderAddError(err error, userID int, w http.ResponseWriter) {
	var orderAlreadyAddedError *services.OrderAlreadyAddedError
	if errors.As(err, &orderAlreadyAddedError) {
		existingOrder := orderAlreadyAddedError.Order
		if existingOrder.CreatedBy == userID {
			http.Error(w, "order already added", http.StatusOK)
			return
		}
		http.Error(w, "order already added by another user", http.StatusConflict)
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("unexpected error in add order handler")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
